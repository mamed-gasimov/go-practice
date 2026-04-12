package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const REQUEST_LIMIT = 1000

type Request struct {
	Payload string
}

type Client interface {
	SendRequest(ctx context.Context, request Request) error
	WithLimiter(ctx context.Context, ch chan Request)
}

type client struct{}

func (c client) SendRequest(ctx context.Context, request Request) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Sending request ", request.Payload)
	return nil
}

var maxConnections = 10

func (c client) WithLimiter(ctx context.Context, ch <-chan Request) {
	wg := sync.WaitGroup{}
	wg.Add(maxConnections)

	for range maxConnections {
		go func() {
			defer wg.Done()
			for req := range ch {
				c.SendRequest(ctx, req)
			}
		}()
	}

	wg.Wait()
}

func main() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, REQUEST_LIMIT)

	for i := 0; i < REQUEST_LIMIT; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.WithLimiter(ctx, generate(requests))
}

func generate(requests []Request) <-chan Request {
	ch := make(chan Request)

	go func() {
		for _, req := range requests {
			ch <- req
		}
		close(ch)
	}()

	return ch
}
