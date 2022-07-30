package transport

import (
	"MyEntryTask/models"
	"MyEntryTask/utils"
	"bufio"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

type Transport struct {
	conn net.Conn
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{
		conn: conn,
	}
}

func (ts *Transport) Send(msg *models.RPCMessage) error {
	w := bufio.NewWriter(ts.conn)

	header, err := proto.Marshal(msg.Header)
	if err != nil {
		utils.Logs.Warn("proto marshal header err :[%v]\n", err)
		return err
	}

	body := msg.Body
	lenHeader := len(header)
	lenBody := len(body)
	msg.HeaderLen = uint64(lenHeader)
	msg.BodyLen = uint64(lenBody)

	buf := make([]byte, 8+lenHeader)
	// put header len
	binary.BigEndian.PutUint64(buf[:8], uint64(lenHeader))
	// put header bytes
	copy(buf[8:], header)
	_, err = w.Write(buf)
	if err != nil {
		utils.Logs.Warn("write buf :[%v] err :[%v]\n", buf, err)
		return err
	}
	_ = w.Flush()

	buf = make([]byte, 8+lenBody)
	// put body len
	binary.BigEndian.PutUint64(buf[:8], uint64(lenBody))
	// put body bytes
	copy(buf[8:], body)
	_, err = w.Write(buf)
	if err != nil {
		utils.Logs.Warn("write buf :[%v] err :[%v]\n", buf, err)
		return err
	}
	_ = w.Flush()

	return nil
}

func (ts *Transport) Receive(msg *models.RPCMessage) (err error) {
	r := bufio.NewReader(ts.conn)
	msg.Header = &models.Header{
		Comm:       &models.CommHeader{},
		ReqHeader:  &models.ReqHeader{},
		RespHeader: &models.RespHeader{},
	}
	// read header len
	buf := make([]byte, 8)
	_, err = r.Read(buf)
	if err != nil {
		if err != io.EOF {
			utils.Logs.Warn("read conn to buf err :[%v]\n", err)
		}
		return err
	}
	msg.HeaderLen = binary.BigEndian.Uint64(buf)
	// read header
	header := make([]byte, msg.HeaderLen)
	_, err = r.Read(header)
	if err != nil {
		if err != io.EOF {
			utils.Logs.Warn("read conn to buf err :[%v]\n", err)
		}
		return err
	}

	if err := proto.Unmarshal(header, msg.Header); err != nil {
		utils.Logs.Warn("proto unmarshal header err :[%v]\n", err)
		return err
	}

	// read body len
	buf = make([]byte, 8)
	_, err = r.Read(buf)
	if err != nil {
		if err != io.EOF {
			utils.Logs.Warn("read conn to buf err :[%v]\n", err)
		}
		return err
	}
	msg.BodyLen = binary.BigEndian.Uint64(buf)
	// read body
	body := make([]byte, msg.BodyLen)
	_, err = r.Read(body)
	if err != nil {
		if err != io.EOF {
			utils.Logs.Warn("read conn to buf err :[%v]\n", err)
		}
		return err
	}
	msg.Body = body
	return nil
}

func (ts *Transport) Close() error {
	return ts.conn.Close()
}
