package geerpc

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

func (m *methodType) NumCalls() uint64 {
	//原子操作加载unint64类型的值
	return atomic.LoadUint64(&m.numCalls)
}

// 两种argv类型：值类型和指针类型，example输出：
// Pointer type: 0x140000a2018, Type: *int
// Value type: 0, Type: int
func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	// reflect相关解释：
	// reflect.New: 返回一个Value类型，该类型包含一个指向类型为typ的新申请的零值
	// reflect.Value:反射包中一个代表Go值的通用容器，可以操作任意类型的值
	// 指针类型和值类型创建实例的方式有细微区别
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}

	return argv
}

// reply必须是指针类型
func (m *methodType) newReplyv() reflect.Value {
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

func newService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Fatal("rpc server: ", s.name, " is not a valid service name")
	}
	s.registerMethods()
	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		//作为rpc调用，包含自身在内，numIn必须为3（自身/导出的方法/返回值指针）;numOut必须为1（error返回值）
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		//(*error)(nil)用于得到一个值为nil的error指针类型;Elem()会直接返回'error'这个类型
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	//returnValues是一个Value类型的切片
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
