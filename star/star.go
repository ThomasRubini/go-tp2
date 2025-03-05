package main

import (
	"math/rand"
	"time"
)

func master(input, output chan int) {
	for {
		// Print all received values
		// (Read from output channel in a non-blocking way)
		moreToReceive := true
		for moreToReceive {
			select {
			case value := <-output:
				println(value)
			default:
				moreToReceive = false
			}
		}

		// Generate number
		randomNumber := rand.Intn(100) // generates a random number between 0 and 99
		input <- randomNumber

		// sleep
		time.Sleep(100 * time.Millisecond)
	}
}

func slave(input, output chan int) {
	for {
		value := <-input
		output <- value * value
	}
}

func main() {
	input := make(chan int)
	output := make(chan int)

	M := 2
	for i := 0; i < M; i++ {
		go slave(input, output)
	}

	master(input, output)
}
