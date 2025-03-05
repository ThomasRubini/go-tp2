package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Special node to reset counter. Does not do work
func resetNode(nodeID int32, ch chan int, wantedNode *int32, resetFunction func(*int) bool) {
	for {
		val := <-ch
		if atomic.LoadInt32(wantedNode) == nodeID {
			atomic.StoreInt32(wantedNode, 0) // reset

			cont := resetFunction(&val) // call callback with final value
			if !cont { // check if we should continue
				break
			}

			ch <- val
		} else {
			// We weren't the right node. Put it back for the correct node.
			ch <- val
		}
	}
}

func node(nodeID int32, ch chan int, wantedNode *int32, workFunction func(int) int) {
	for {
		val := <-ch
		if atomic.LoadInt32(wantedNode) == nodeID {
			fmt.Println("Goroutine", nodeID, "received:", val)

			// Signify that the next node can take it
			atomic.AddInt32(wantedNode, 1)

			// Make work + pass to next
			ch <- workFunction(val)
		} else {
			// We weren't the right node. Put it back for the correct node.
			ch <- val
		}
	}
}

func work(value int) int {
	return value + 5
}

func main() {
	P := 5
	ch := make(chan int)
	var wantedNode int32 = 0

	for i := 0; i < P; i++ {
		go node(int32(i), ch, &wantedNode, work)
	}

	rounds := 0
	K := 10
	go resetNode(int32(P), ch, &wantedNode, func(i *int) bool {
		rounds++
		if rounds >= K {
			fmt.Printf("Everything is finished ! Final value: %v\n", *i)
			return false
		} else {
			return true
		}
	})

	ch <- 0
	time.Sleep(100 * time.Second)
}
