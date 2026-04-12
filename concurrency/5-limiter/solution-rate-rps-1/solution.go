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
	WithLimiter(ctx context.Context, requests []Request)
}

type client struct{}

func (c client) SendRequest(ctx context.Context, request Request) error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Sending request ", request.Payload)
	return nil
}

var rps = 100

func (c client) WithLimiter(ctx context.Context, requests []Request) {
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	wg := sync.WaitGroup{}
	wg.Add(len(requests))

	for _, req := range requests {
		<-ticker.C
		go func() {
			defer wg.Done()
			c.SendRequest(ctx, req)
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
	c.WithLimiter(ctx, requests)
}
