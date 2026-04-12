package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const REQUEST_LIMIT = 10

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

var maxGoroutines = 100

func (c client) WithLimiter(ctx context.Context, requests []Request) {
	// создали хранилище с токенами
	tokens := make(chan struct{}, maxGoroutines)

	// заполнили хрпнилище полностью
	go func() {
		for range maxGoroutines {
			tokens <- struct{}{}
		}
	}()

	for _, req := range requests {
		// каждая горутина до старта забирает один токен
		// если возьмет больше, чем размер хранилища, то блокируется тут
		<-tokens
		go func() {
			// затеи после завершения горутина возвращает токен в хранилище
			defer func() {
				tokens <- struct{}{}
			}()

			c.SendRequest(ctx, req)
		}()
	}

	// для синхронизации считываем
	for range maxGoroutines {
		<-tokens
	}
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
