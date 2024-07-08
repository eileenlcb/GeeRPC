package main

import "time"

func fetchData(result chan<- string) {
	time.Sleep(3 * time.Second)
	result <- "fetch data success"
}

func main() {
	resultChan := make(chan string)
	go fetchData(resultChan)
	select {
	case <-time.After(4 * time.Second):
		println("timeout")
	case result := <-resultChan:
		println("success:", result)
	}
}
