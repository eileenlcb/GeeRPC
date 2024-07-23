package main

import "time"

type result struct {
	value int
	err   error
}

func cal(n int, ch chan result) {
	time.Sleep(time.Second * 2)
	var r result
	r.value = n * n
	ch <- r
}

func main() {
	ch := make(chan result)
	go cal(10, ch)

	select {
	case res := <-ch:
		if res.err != nil {
			println(res.err)
		} else {
			println(res.value)
		}
	case <-time.After(time.Second * 3):
		println("timeout")
	}
}
