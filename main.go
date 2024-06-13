package main

import (
	"fmt"
	"reflect"
)

type methodType struct {
	ArgType reflect.Type
}

// 添加 newArgv 方法
func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	// 两种argv类型：值类型和指针类型
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

func main() {
	var m methodType

	// 指针类型
	m = methodType{ArgType: reflect.TypeOf((*int)(nil))}
	argvPtr := m.newArgv()
	fmt.Printf("Pointer type: %v, Type: %v\n", argvPtr, argvPtr.Type())

	// 非指针类型
	m = methodType{ArgType: reflect.TypeOf(int(0))}
	argvVal := m.newArgv()
	fmt.Printf("Value type: %v, Type: %v\n", argvVal, argvVal.Type())
}
