package geerpc

import (
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

func (m *methodType) NumCalls() uint64 {
	//原子操作加载unint64类型的值
	return atomic.LoadUint64(&m.numCalls)
}

// reflect相关解释：
// reflect.New: 返回一个Value类型，该类型包含一个指向类型为typ的新申请的零值
// reflect.Value:反射包中一个代表Go值的通用容器，可以操作任意类型的值

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	//两种argv类型：值类型和指针类型，example输出：
	// Pointer type: 0x140000a2018, Type: *int
	// Value type: 0, Type: int
	// 指针类型和值类型创建实例的方式有细微区别
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}

	return argv
}

func (m *methodType) newReplyv() reflect.Value {
	// reply必须是指针类型
	// 指向m.ReplyType.elem()实例的指针
	replyv := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}
