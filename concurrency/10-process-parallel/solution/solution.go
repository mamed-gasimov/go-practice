package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type outVal struct {
	val int
	err error
}

var errTimeout = errors.New("timed out")

func processData(ctx context.Context, val int) chan outVal {
	ch := make(chan struct{})
	out := make(chan outVal)

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		close(ch)
	}()

	go func() {
		select {
		case <-ch:
			out <- outVal{
				val: val * 2,
				err: nil,
			}
		case <-ctx.Done():
			out <- outVal{
				val: 0,
				err: errTimeout,
			}
		}
	}()

	return out
}

func main() {
	in := make(chan int)
	out := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		defer close(in)
		for i := range 10 {
			select {
			case in <- i + 1:
			case <-ctx.Done():
				return
			}
		}
	}()

	now := time.Now()
	processParallel(ctx, in, out, 5)

	for val := range out {
		fmt.Println(val)
	}
	fmt.Println(time.Since(now))
}

// операция должна выполняться не более пяти секунд
func processParallel(ctx context.Context, in <-chan int, out chan<- int, numWorkers int) {
	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)

	for range numWorkers {
		go worker(ctx, in, out, wg)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
}

func worker(ctx context.Context, in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case val, ok := <-in:
			if !ok {
				return
			}
			select {
			case ov := <-processData(ctx, val):
				if ov.err != nil {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- ov.val:
				}
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}

}
