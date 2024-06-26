package geerpc

import (
	"fmt"
	"reflect"
	"testing"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

// it's not a exported Method
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func TestNewService(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	_assert(len(s.method) == 1, "wrong service Method, expect 1, but got %d", len(s.method))
	mType := s.method["Sum"]
	_assert(mType != nil, "wrong Method, Sum shouldn't nil")
}

func TestMethodType_Call(t *testing.T) {
	var foo Foo
	//&geerpc.service{name:"Foo", typ:(*reflect.rtype)(0x102667f20), 
	//rcvr:reflect.Value{typ_:(*abi.Type)(0x102667f20), ptr:(unsafe.Pointer)(0x140001a2310), flag:0x16}, 
	//method:map[string]*geerpc.methodType{"Sum":(*geerpc.methodType)(0x140001e0080)}}
	s := newService(&foo)
	mType := s.method["Sum"]
	fmt.Println("mType:", mType)
	argv := mType.newArgv()
	replyv := mType.newReplyv()
	fmt.Println("argv:", argv)
	fmt.Println("replyv:", replyv)
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 3}))
	err := s.call(mType, argv, replyv)
	_assert(err == nil && *replyv.Interface().(*int) == 4 && mType.NumCalls() == 1, "failed to call Foo.Sum")
}

func TestFindService(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	fmt.Printf("%#v\n", s)
	serviceMethod := "Foo.Sum"
	svc, mType, _ := DefaultServer.findService(serviceMethod)
	fmt.Println("#######")
	fmt.Println("svc:", svc)
	fmt.Println("mType:", mType)
}
