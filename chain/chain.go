package main

import "time"

func main() {
	ch := make(chan int)
	N := 10
	for i := 0; i < N; i++ {
		go func() {
			value := <-ch
			value++

			println(value)

			if value == N-1 {
				println("FIN")
				println(value)
			} else {
				ch <- value
			}
		}()
	}

	ch <- 0
	time.Sleep(1 * time.Second) // Prevent main from exiting early
}
