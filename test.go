package main

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type service struct {
	rcvr reflect.Value
}

type methodType struct {
	method   reflect.Method
	numCalls uint64
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

// 一个测试方法
func (s *service) ExampleMethod(arg int, reply *int) error {
	*reply = arg + 1
	return nil
}

func main() {
	srv := &service{}
	srv.rcvr = reflect.ValueOf(srv)

	exampleMethod, _ := reflect.TypeOf(srv).MethodByName("ExampleMethod")
	m := &methodType{method: exampleMethod}

	argv := reflect.ValueOf(42)
	// 创建指针类型的 reflect.Value
	replyv := reflect.New(reflect.TypeOf(0)).Elem().Addr()

	err := srv.call(m, argv, replyv)
	fmt.Println("Error:", err)
	// 获取指针指向的值
	fmt.Println("Reply:", replyv.Elem().Int())
}
