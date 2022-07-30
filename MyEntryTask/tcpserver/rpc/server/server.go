package server

import (
	"MyEntryTask/models"
	"MyEntryTask/tcpserver/rpc/transport"
	"MyEntryTask/utils"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"
)

type RpcSerConf struct {
	addr string
	port string
}

type RpcServer struct {
	serviceMap sync.Map
	conf       *RpcSerConf
	longconn   bool
}

func NewRpcServerConf(addr, port string) *RpcSerConf {
	return &RpcSerConf{
		addr: addr,
		port: port,
	}
}

func NewRpcServer(conf *RpcSerConf) *RpcServer {
	return &RpcServer{
		conf:     conf,
		longconn: true,
	}
}

func (svr *RpcServer) Run() {
	addr := svr.conf.addr + ":" + svr.conf.port
	//utils.Logs.Info("rpc run at addr: [%s]\n", addr)

	lis, err := net.Listen("tcp", addr)
	utils.Logs.Info("rpc svr run at :[%s]\n", addr)
	if err != nil {
		utils.Logs.Error("net listen failed, addr :[%s], err :[%v]\n", addr, err)
		return
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			utils.Logs.Error("rpc server accept err: [%v]\n", err)
			return
		}

		go svr.HandleConn(conn)
	}
}

// HandleConn 处理rpc请求
func (svr *RpcServer) HandleConn(conn net.Conn) {
	ts := transport.NewTransport(conn)
	defer ts.Close()
	defer utils.RecoverPanic()

	for {
		req := &models.RPCMessage{}
		resp := &models.RPCMessage{
			Header: &models.Header{
				Comm:       &models.CommHeader{},
				ReqHeader:  &models.ReqHeader{},
				RespHeader: &models.RespHeader{},
			},
		}

		err := ts.Receive(req)
		//utils.Logs.Debug("svr recv :[%+v]\n", *req)
		if err == io.EOF {
			return
		}
		if err != nil {
			utils.Logs.Error("rpc receive err: [%v]\n", err)
			utils.SetRespCode(resp, uint32(utils.CodeTCPRpcTransportErr))
			_ = ts.Send(resp)
			return
		}
		resp.Header = req.Header

		svc, mtype, err := svr.findSvcMethod(req.Header.Comm.Method)
		if err != nil {
			utils.Logs.Warn("find svr method err :[%v]\n", err)
			utils.SetRespCode(resp, uint32(utils.CodeTCPRpcNotFindSvcOrMethod))
			_ = ts.Send(resp)
			return
		}

		// 服务，方法，请求参数，响应参数
		err = svr.callMethod(svc, mtype, req, resp)
		if err != nil {
			utils.Logs.Warn("call method err :[%v]\n", err)
			utils.SetRespCode(resp, uint32(utils.CodeTCPRpcServiceErr))
			_ = ts.Send(resp)
			return
		}

		// send resp
		if err := ts.Send(resp); err != nil {
			utils.Logs.Warn("rpc send err: [%v]\n", err)
			return
		}

		if !svr.longconn {
			return
		}
	}

}

// findSvrMethod 找对应的服务和方法
func (svr *RpcServer) findSvcMethod(svcMethod string) (svc *models.Service, mType *models.MethodType, err error) {
	dot := strings.LastIndex(svcMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + svcMethod)
		return
	}
	svcName, methodName := svcMethod[:dot], svcMethod[dot+1:]
	servicei, ok := svr.serviceMap.Load(svcName)
	if !ok {
		err = errors.New("rpc server: can't find service " + svcName)
		return
	}

	svc = servicei.(*models.Service)
	mType = svc.Method[methodName]
	if mType == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}

	return
}

// callMethod 方法调用
func (svr *RpcServer) callMethod(svc *models.Service, mthed *models.MethodType, req, resp *models.RPCMessage) (err error) {
	f := mthed.Method.Func
	argv := reflect.New(mthed.ArgType.Elem())
	replyv := reflect.New(mthed.ReplyType.Elem())

	if err := proto.Unmarshal(req.Body, argv.Interface().(proto.Message)); err != nil {
		utils.Logs.Warn("proto unmarshal err :[%v]\n", err)
		return err
	}

	returnVal := f.Call([]reflect.Value{svc.Rcvr, reflect.ValueOf(req.Header), argv, replyv})
	retErr := returnVal[0].Interface()
	if retErr != nil {
		err = retErr.(error)
		utils.Logs.Warn("rpc call func [%v] err :[%v]\n", mthed.Method, err)
		return err
	}

	resp.Body, err = proto.Marshal(replyv.Interface().(proto.Message))
	if err != nil {
		utils.Logs.Warn("proto marshal err :[%v]\n", err)
		return err
	}

	// 设置resp header
	resp.Header = req.Header
	return nil
}

func (svr *RpcServer) NewService() *models.Service {
	return &models.Service{}
}

// Register 服务注册
func (svr *RpcServer) Register(rcvr interface{}) (err error) {
	svc := svr.NewService()
	svc.Typ = reflect.TypeOf(rcvr)
	svc.Rcvr = reflect.ValueOf(rcvr)
	svc.Name = reflect.Indirect(svc.Rcvr).Type().Name()

	// register method
	svc.Method = make(map[string]*models.MethodType)

	for i := 0; i < svc.Typ.NumMethod(); i++ {
		method := svc.Typ.Method(i)
		mtype := method.Type

		// 校验service api
		if mtype.NumIn() != 4 || mtype.NumOut() != 1 || mtype.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			utils.Logs.Warn("method :[%s] invalid\n", method.Name)
			continue
		}

		argsType := mtype.In(2)
		replyType := mtype.In(3)
		svc.Method[method.Name] = &models.MethodType{
			Method:    method,
			ArgType:   argsType,
			ReplyType: replyType,
		}
		//utils.Logs.Info("rpc server: register %s.%s\n", svc.Name, method.Name)
	}

	// svc存在svr的sync.map
	_, ok := svr.serviceMap.LoadOrStore(svc.Name, svc)
	if ok {
		utils.Logs.Warn("rpc server: service already defined:%s\n", svc.Name)
		return errors.New("rpc: service already defined: " + svc.Name)
	}
	return nil
}
