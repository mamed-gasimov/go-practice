package main

import (
	"context"
	"sync"
)

func fanIn(ctx context.Context, chans []chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(chans))

	worker := func(ch chan int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-ch:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- value:
				}
			}
		}
	}

	for _, ch := range chans {
		go worker(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {}
