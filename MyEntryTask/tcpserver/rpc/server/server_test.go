package server

import (
	"MyEntryTask/models"
	"MyEntryTask/tcpserver/rpc/client"
	"MyEntryTask/utils"
	"context"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestRpcSvr(t *testing.T) {
	conf := &RpcSerConf{
		addr: "127.0.0.1",
		port: "10002",
	}
	rpc := NewRpcServer(conf)

	// service
	serviceTestCall := &TestCall{}
	reqbody := &models.TestRpcCallReq{
		Username: "zhangjunhao",
		Password: "@shopee",
	}

	// service reply
	respBody := &models.TestRpcCallResp{}

	err := rpc.Register(serviceTestCall)
	if err != nil {
		t.Errorf("rpc register err :[%v]\n", err)
		return
	}

	svc, mtype, err := rpc.findSvcMethod("TestCall.Link")
	if err != nil {
		t.Errorf("rpc find svc method err :[%v]\n", err)
		return
	}

	req := &models.RPCMessage{}
	resp := &models.RPCMessage{}
	req.Header = &models.Header{
		Comm: &models.CommHeader{
			CallSeq: uint64(1),
			Method:  "TestCall.Link",
		},
		ReqHeader: &models.ReqHeader{},
		RespHeader: &models.RespHeader{
			Code: utils.CodeSucc,
		},
	}

	req.Body, _ = proto.Marshal(reqbody)
	err = rpc.callMethod(svc, mtype, req, resp)
	if err != nil {
		t.Errorf("rpc call method err :[%v]\n", err)
		return
	}

	if err := proto.Unmarshal(resp.Body, respBody); err != nil {
		t.Errorf("proto unmarshal err :[%v]\n", err)
		return
	}
	if respBody.Result != "zhangjunhao@shopee" {
		t.Errorf("resp body err return value is :[%s]\n", respBody.Result)
		return
	}
	if resp.Header.RespHeader.Code != 1 {
		t.Errorf("resp code uncorrect return value is :[%d]\n", resp.Header.RespHeader.Code)
		return
	}
}

type TestCall struct {
}

func (t *TestCall) Link(header *models.Header, arg *models.TestRpcCallReq, reply *models.TestRpcCallResp) error {
	utils.SetHeaderCode(header, utils.CodeSucc)
	reply.Result = arg.Username + arg.Password
	return nil
}

func TestRpcDemo(t *testing.T) {
	conf := &RpcSerConf{
		addr: "127.0.0.1",
		port: "8888",
	}
	rpc := NewRpcServer(conf)
	go rpc.Run()
	// service
	serviceTestCall := &TestCall{}
	args := &models.TestRpcCallReq{
		Username: "zhangjunhao",
		Password: "@shopee",
	}

	// service reply
	reply := &models.TestRpcCallResp{}

	err := rpc.Register(serviceTestCall)
	if err != nil {
		t.Errorf("rpc register err :[%v]\n", err)
		return
	}

	//svc, mtype, err := rpc.findSvcMethod("TestCall.Link")
	//if err != nil {
	//	t.Errorf("rpc find svc method err :[%v]\n", err)
	//	return
	//}
	//
	//req := &models.RPCMessage{}
	//resp := &models.RPCMessage{}
	//req.Header = &models.Header{
	//	Comm: &models.CommHeader{
	//		CallSeq: uint64(1),
	//		Method:  "TestCall.Link",
	//	},
	//	ReqHeader: &models.ReqHeader{},
	//	RespHeader: &models.RespHeader{
	//		Code: utils.CodeSucc,
	//	},
	//}
	poolConf := &client.CliPoolConf{
		Size:   2000,
		MaxCap: 3000,
	}

	pool, err := client.NewClientPool(poolConf)
	if err != nil {
		t.Errorf("new pool err :[%v]\n", err)
		return
	}

	cli, err := pool.Get()
	if err != nil {
		t.Errorf("pool get err :[%v]\n", err)
		return
	}

	ctx, errCtx := context.WithTimeout(context.Background(), 5*time.Second)
	if errCtx != nil {
		t.Errorf("context err:[%v]\n", err)
		return
	}
	code, err := cli.CallRpc(ctx, "TestCall.Link", args, reply)
	if err != nil {
		t.Errorf("rpc call err, code :[%d] err:[%v]\n", code, err)
		return
	}

	if reply.Result != "zhangjunhao@shopee" {
		t.Errorf("resp body err return value is :[%s]\n", reply.Result)
		return
	}
	if code != utils.CodeSucc {
		t.Errorf("resp code uncorrect return value is :[%d]\n", code)
		return
	}
}
