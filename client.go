package geerpc

type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         string
	Done          chan *Call
}

func (call *Call) done() {
	call.Done <- call
}
