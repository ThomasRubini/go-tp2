package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Hello from goroutine"
	}()

	msg := <-ch
	fmt.Println(msg)
}
