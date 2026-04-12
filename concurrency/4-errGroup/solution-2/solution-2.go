package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
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

	egroup, ectx := errgroup.WithContext(ctx)
	// limit of active goroutines
	egroup.SetLimit(100)

	for _, u := range users {
		egroup.Go(func() error {
			name, err := fetch(ectx, u)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			names[name] = names[name] + 1
			return nil
		})
	}

	if err := egroup.Wait(); err != nil {
		return nil, err
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
