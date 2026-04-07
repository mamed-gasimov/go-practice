package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Backend struct {
	Name  string
	Delay time.Duration
}

func (w *Backend) DoRequest(ctx context.Context) string {
	select {
	case <-ctx.Done():
		return ""
	case <-time.After(w.Delay):
	}

	return fmt.Sprintf("Response of backend %s", w.Name)
}

func main() {
	backends := []*Backend{
		{Name: "NYSE", Delay: 4 * time.Second},
		{Name: "MOEX", Delay: 1 * time.Second},
		{Name: "LSE", Delay: 2 * time.Second},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var results []string
	for name := range getFirstName(ctx, backends) {
		results = append(results, name)
	}

	<-ctx.Done()
	fmt.Println(results)
}

func getFirstName(ctx context.Context, backends []*Backend) <-chan string {
	names := make(chan string)
	var wg sync.WaitGroup

	for _, backend := range backends {
		wg.Add(1)
		go func() {
			names <- backend.DoRequest(ctx)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(names)
	}()

	return names
}
