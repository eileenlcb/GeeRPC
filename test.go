package main

import (
	"fmt"
	"reflect"
)

func add(a, b int) int {
	return a + b
}

func main() {
	f := reflect.ValueOf(add)

	args := []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)}
	result := f.Call(args)

	sum := result[0].Int()
	fmt.Println(sum)
}
