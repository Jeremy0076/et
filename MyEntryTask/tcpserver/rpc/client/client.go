package client

import (
	"MyEntryTask/models"
	"MyEntryTask/tcpserver/rpc/transport"
	"MyEntryTask/utils"
	"context"
	"errors"
	"google.golang.org/protobuf/proto"
	"io"
	"math/rand"
	"net"
	"time"
)

const RpcHost = "127.0.0.1:8889"

type RPCClient struct {
	longconn bool
	conn     net.Conn
	reqId    uint32
}

// Call 调用
type Call struct {
	Done chan *Call // 在调用结束时激活
}

func NewRpcClient() *RPCClient {
	return &RPCClient{
		longconn: true,
		reqId:    uint32(rand.Int()),
	}
}

func (cli *RPCClient) Close() {
	_ = cli.conn.Close()
}

func (cli *RPCClient) DialRPC(call chan interface{}, network string, addr string) {
	err := cli.Dial(call, network, addr)
	if err != nil {
		utils.Logs.Warn("dial rpc err :[%v]\n", err)
		return
	}
}

func (cli *RPCClient) Dial(call chan interface{}, network string, addr string) (err error) {
	cli.conn, err = net.Dial(network, RpcHost)
	if err != nil {
		utils.Logs.Warn("dial tcp addr :[%s] err :[%v]\n", addr, err)
		return err
	}
	call <- 1
	return nil
}

// CallRpc call rpc
func (cli *RPCClient) CallRpc(ctx context.Context, svcMethod string, args, reply proto.Message) (code uint32, err error) {
	if cli.conn == nil {
		callDial := make(chan interface{}, 1)
		go cli.DialRPC(callDial, "tcp", RpcHost)
		// 连接超时处理
		select {
		case <-time.After(2 * time.Second):
			return utils.CodeTCPRpcTimeout, nil
		case _ = <-callDial:
			//utils.Logs.Info("rpc call dial addr:[%s]\n", RpcHost)
		}
	}
	respMsg := &models.RPCMessage{}
	reqMsg := &models.RPCMessage{
		Header: &models.Header{
			Comm:       &models.CommHeader{},
			ReqHeader:  &models.ReqHeader{},
			RespHeader: &models.RespHeader{},
		},
	}
	if cli.conn != nil && cli.longconn {
		// 封装请求

		reqMsg.Header.Comm.Method = svcMethod
		reqMsg.Header.Comm.CallSeq = uint64(cli.reqId)
		cli.reqId++
		reqMsg.Body, err = proto.Marshal(args)
		if err != nil {
			utils.Logs.Warn("proto marshal err:[%v]\n", err)
			return utils.CodeTCPInternelErr, err
		}

		callGo := make(chan interface{}, 1)

		go cli.GoRpc(callGo, reqMsg, respMsg)

		// 超时处理
		select {
		case <-ctx.Done():
			return utils.CodeTCPRpcTimeout, errors.New("rpc client: call failed: " + ctx.Err().Error())
		case _ = <-callGo:
		}

		if respMsg.Header.RespHeader.Code != utils.CodeSucc {
			utils.Logs.Warn("rpc call code err: [%d]\n", respMsg.Header.RespHeader.Code)
			return respMsg.Header.RespHeader.Code, err
		}
		// 封装响应
		err = proto.Unmarshal(respMsg.Body, reply)
		if err != nil {
			utils.Logs.Warn("proto unmarshal err:[%v]\n", err)
			return utils.CodeTCPInternelErr, err
		}
		return utils.CodeSucc, nil
	}

	retry := 3
	for ; retry > 0; retry-- {
		cli.conn, err = net.Dial(kNetwork, RpcHost)
		if err != nil {
			utils.Logs.Warn("dial failed addr :[%s] err :[%v]\n", RpcHost, err)
			continue
		}
		break
	}

	if retry == 0 {
		utils.Logs.Warn("dial failed 3 times addr :[%s] err :[%v]\n", RpcHost, err)
		return utils.CodeTCPRpcServiceErr, err
	}

	callGo := make(chan interface{}, 1)

	go cli.GoRpc(callGo, reqMsg, respMsg)

	// 超时处理
	select {
	case <-ctx.Done():
		return utils.CodeTCPRpcTimeout, errors.New("rpc client: call failed: " + ctx.Err().Error())
	case _ = <-callGo:
	}

	if respMsg.Header.RespHeader.Code != utils.CodeSucc {
		utils.Logs.Warn("rpc call code err: [%d]\n", respMsg.Header.RespHeader.Code)
		return respMsg.Header.RespHeader.Code, err
	}
	// 封装响应
	err = proto.Unmarshal(respMsg.Body, reply)
	if err != nil {
		utils.Logs.Warn("proto unmarshal err:[%v]\n", err)
		return utils.CodeTCPInternelErr, err
	}

	if !cli.longconn {
		cli.Close()
		cli.conn = nil
	}
	return utils.CodeSucc, nil
}

func (cli *RPCClient) goRpc(call chan interface{}, req, resp *models.RPCMessage) (err error) {
	ts := transport.NewTransport(cli.conn)
	err = ts.Send(req)
	if err != nil {
		utils.Logs.Warn("rpc client send req msg err:[%v]\n", err)
		return err
	}
	//utils.Logs.Debug("cli send :[%+v]\n", *req)

	err = ts.Receive(resp)
	if err != nil && err != io.EOF {
		utils.Logs.Warn("rpc client receive req msg err:[%v]\n", err)
		return err
	}

	//utils.Logs.Debug("cli recv :[%+v]\n", *resp)
	call <- 1
	return nil
}

func (cli *RPCClient) GoRpc(call chan interface{}, req, resp *models.RPCMessage) {
	err := cli.goRpc(call, req, resp)
	if err != nil {
		utils.Logs.Warn("go rpc err :[%v]\n", err)
		return
	}
}
