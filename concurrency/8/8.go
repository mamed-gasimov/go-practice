package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randomWait() int {
	workSeconds := rand.Intn(6)

	time.Sleep(time.Duration(workSeconds) * time.Second)

	return workSeconds
}

func solution1() int {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	totalWorkSeconds := 0

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			seconds := randomWait()

			mu.Lock()
			totalWorkSeconds += seconds
			mu.Unlock()
		}()
	}

	wg.Wait()
	return totalWorkSeconds
}

func solution2() int {
	ch := make(chan int)
	totalWorkSeconds := 0

	for range 100 {
		go func() {
			ch <- randomWait()
		}()
	}

	for range 100 {
		totalWorkSeconds += <-ch
	}

	return totalWorkSeconds
}

func main() {
	start := time.Now()
	totalWorkSeconds := solution1()
	mainSeconds := time.Since(start)
	fmt.Println("main 1: ", mainSeconds)
	fmt.Println("total 1: ", totalWorkSeconds)

	start2 := time.Now()
	totalWorkSeconds2 := solution2()
	mainSeconds2 := time.Since(start2)
	fmt.Println("main 2: ", mainSeconds2)
	fmt.Println("total 2: ", totalWorkSeconds2)
}
