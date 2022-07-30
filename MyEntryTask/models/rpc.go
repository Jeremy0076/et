package models

type RPCMessage struct {
	HeaderLen uint64
	Header    *Header
	BodyLen   uint64
	Body      []byte
}
