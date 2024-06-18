package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

type GobCodec struct {
	dec *gob.Decoder
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func main() {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode("test test")

	dec := gob.NewDecoder(&buf)
	cc := GobCodec{dec: dec}

	type request struct {
		h             interface{}
		argv, replayv reflect.Value
	}
	req := &request{}

	req.argv = reflect.New(reflect.TypeOf(""))
	_ = cc.ReadBody(req.argv.Interface())
	fmt.Print("argv:", req.argv.Elem().String(), "\n")
}
