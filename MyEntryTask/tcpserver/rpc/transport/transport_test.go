package transport

import (
	"MyEntryTask/models"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
	"testing"
)

func TestTransport(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	bodyValue := &models.TestTransport_User{
		Username: "zhangjunhao",
		Password: "shopee",
	}

	rpcMsg := &models.RPCMessage{
		Header: &models.Header{
			Comm: &models.CommHeader{
				CallSeq: uint64(1),
				Method:  "test",
			},
		},
	}
	rpcMsg.Body, _ = proto.Marshal(bodyValue)

	addr := "127.0.0.1:10009"

	// receiver
	go func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic err :[%v]\n", err)
			}
		}()
		conn, _ := net.Dial("tcp", addr)
		ts := NewTransport(conn)

		defer wg.Done()
		err := ts.Send(rpcMsg)
		if err != nil {
			t.Errorf("send msg err")
		}
	}()

	// sender
	go func() {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Panic err :[%v]\n", err)
			}
		}()
		defer wg.Done()
		lis, _ := net.Listen("tcp", addr)
		defer lis.Close()
		conn, _ := lis.Accept()
		defer conn.Close()

		ts := NewTransport(conn)
		msg := &models.RPCMessage{}

		err := ts.Receive(msg)
		if err != nil {
			t.Errorf("receive msg err")
		}

		testBody := &models.TestTransport_User{}
		_ = proto.Unmarshal(msg.Body, testBody)
		if testBody.Username != bodyValue.Username {
			t.Errorf("body value uncorrect")
		}

		if testBody.Password != bodyValue.Password {
			t.Errorf("body value uncorrect")
		}

		if msg.Header.Comm.CallSeq != rpcMsg.Header.Comm.CallSeq {
			t.Errorf("header callseq uncorrect")
		}

		if msg.Header.Comm.Method != rpcMsg.Header.Comm.Method {
			t.Errorf("header method uncorrect")
		}
		fmt.Println("4")
	}()
	wg.Wait()
}
