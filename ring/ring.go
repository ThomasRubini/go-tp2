package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Special node to reset counter. Does not do work
func (r *Ring) resetNode(nodeID int32, wantedNode *int32, resetFunction func(*int) bool) {
	for {
		val := <-r.dataChannel
		if atomic.LoadInt32(wantedNode) == nodeID {
			atomic.StoreInt32(wantedNode, 0) // reset

			cont := r.FinishedRound(&val) // call callback with final value
			if !cont {                    // check if we should continue
				break
			}

			r.dataChannel <- val
		} else {
			// We weren't the right node. Put it back for the correct node.
			r.dataChannel <- val
		}
	}
}

func (r *Ring) workWrapper(nodeID int32, wantedNode *int32, workFunction func(int) int) {
	for {
		val := <-r.dataChannel
		if atomic.LoadInt32(wantedNode) == nodeID {
			fmt.Println("Goroutine", nodeID, "received:", val)

			// Signal that the next node can take it
			atomic.AddInt32(wantedNode, 1)

			// Make work + pass to next
			r.dataChannel <- r.Work(val)
		} else {
			// We weren't the right node. Put it back for the correct node.
			r.dataChannel <- val
		}
	}
}

type Ring struct {
	dataChannel   chan int
	Work          func(int) int
	FinishedRound func(*int) bool
}

func NewRing(nodes int) *Ring {
	r := Ring{}
	r.dataChannel = make(chan int)
	var wantedNode int32 = 0

	// Create nodes
	for i := 0; i < nodes; i++ {
		go r.workWrapper(int32(i), &wantedNode, r.Work)
	}

	// create reset node
	go r.resetNode(int32(nodes), &wantedNode, r.FinishedRound)

	return &r
}

func (r Ring) Start(value int) {
	r.dataChannel <- value
}

func main() {
	r := NewRing(10)
	K := 5
	var rounds int

	r.Work = func(i int) int {
		return i + 5
	}

	var wg sync.WaitGroup
	wg.Add(1)
	r.FinishedRound = func(i *int) bool {
		rounds++
		if rounds >= K {
			fmt.Printf("Everything is finished ! Final value: %v\n", *i)
			wg.Done()
			return false
		} else {
			return true
		}
	}

	r.Start(0)
	wg.Wait()
}
