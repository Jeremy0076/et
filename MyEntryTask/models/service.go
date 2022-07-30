package models

import "reflect"

type MethodType struct {
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

type Service struct {
	Name   string
	Typ    reflect.Type
	Rcvr   reflect.Value
	Method map[string]*MethodType
}
