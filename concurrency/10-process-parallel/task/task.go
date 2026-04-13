package main

import (
	"fmt"
	"math/rand"
	"time"
)

func processData(val int) int {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	return val * 2
}

func main() {
	in := make(chan int)
	out := make(chan int)

	go func() {
		for i := range 100 {
			in <- i
		}
		close(in)
	}()

	now := time.Now()
	processParallel(in, out, 5)

	for val := range out {
		fmt.Println(val)
	}
	fmt.Println(time.Since(now))
}

// операция должна выполняться не более пяти секунд
func processParallel(in <-chan int, out <-chan int, numWorkers int) {}
