package main

import (
	"fmt"
	"sync"
	"time"
)

var numWorkers = 5
var numJobs = 30

func producer(jobs []int) <-chan int {
	jobCh := make(chan int)
	go func() {
		for _, job := range jobs {
			jobCh <- job
		}
		close(jobCh)
	}()
	return jobCh
}

func worker(wg *sync.WaitGroup, jobCh <-chan int, w int) {
	defer wg.Done()
	for j := range jobCh {
		fmt.Println("Worker ", w, "started job ", j)
		time.Sleep(time.Second)
		fmt.Println("Worker ", w, "finished job ", j)
	}
}

func main() {
	jobs := make([]int, 0, numJobs)
	for i := range numJobs {
		jobs = append(jobs, i)
	}

	jobCh := producer(jobs)

	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(wg, jobCh, w)
	}

	wg.Wait()
}
