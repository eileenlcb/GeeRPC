package geerpc

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
	//拿到methodType中的方法引用
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
	// rcvr是一个reflect.Value类型，包含了一个指向srv的指针
	srv.rcvr = reflect.ValueOf(srv)

	// reflect.TypeOf(srv)拿到*service类型的reflect.Type
	exampleMethod, _ := reflect.TypeOf(srv).MethodByName("ExampleMethod")
	m := &methodType{method: exampleMethod}

	// 创建一个int类型的 reflect.Value
	argv := reflect.ValueOf(42)
	// reflect.New(reflect.TypeOf(0))创建*int类型的 reflect.Value;Elem()获取指针指向的值，addr()获取指针
	// 发现这行简化为：replyv := reflect.New(reflect.TypeOf(0))也可以
	replyv := reflect.New(reflect.TypeOf(0)).Elem().Addr()

	err := srv.call(m, argv, replyv)
	fmt.Println("Error:", err)
	// 获取指针指向的值
	fmt.Println("Reply:", replyv.Elem().Int())
}
