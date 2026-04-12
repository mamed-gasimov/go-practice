package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type User struct {
	Name string
}

func fetch(_ context.Context, user User) (string, error) {
	time.Sleep(time.Millisecond * 10)
	return user.Name, nil
}

func process(ctx context.Context, users []User) (map[string]int, error) {
	names := make(map[string]int, 0)
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(users))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var commonError error

	for _, u := range users {
		go func() {
			defer wg.Done()

			name, err := fetch(ctx, u)
			if err != nil {
				sync.OnceFunc(func() {
					cancel()
					commonError = err
				})
			}

			mu.Lock()
			defer mu.Unlock()
			names[name] = names[name] + 1
		}()
	}

	wg.Wait()
	if commonError != nil {
		return nil, commonError
	}
	return names, nil
}

func main() {
	names := []User{
		{"Ann"},
		{"Bob"},
		{"Cindy"},
		{"Bob"},
	}

	ctx := context.Background()
	start := time.Now()
	res, err := process(ctx, names)
	if err != nil {
		fmt.Println("An error occured: ", err.Error())
	}

	fmt.Printf("Time passed %v\n", time.Since(start))
	fmt.Println(res)
}
