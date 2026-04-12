package main

import (
	"context"
	"fmt"
	"strconv"
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

func (c client) WithLimiter(ctx context.Context, requests []Request) {}

func main() {
	ctx := context.Background()
	c := client{}
	requests := make([]Request, REQUEST_LIMIT)

	for i := 0; i < REQUEST_LIMIT; i++ {
		requests[i] = Request{Payload: strconv.Itoa(i)}
	}
	c.WithLimiter(ctx, requests)
}
