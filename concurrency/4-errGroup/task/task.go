package main

import (
	"context"
	"fmt"
	"time"
)

// it takes 40ms
// 1. make requests concurrent
// 2. stop all requests one time

type User struct {
	Name string
}

func fetch(_ context.Context, user User) (string, error) {
	time.Sleep(time.Millisecond * 10)
	return user.Name, nil
}

func process(ctx context.Context, users []User) (map[string]int, error) {
	names := make(map[string]int, 0)

	for _, u := range users {
		name, err := fetch(ctx, u)
		if err != nil {
		}

		names[name] = names[name] + 1
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
