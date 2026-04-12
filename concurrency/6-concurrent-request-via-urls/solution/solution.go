package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://amazon.com",
		"https://youtube.com",
	}

	fmt.Println(process(urls))
}

var client http.Client

// реализовать паралельные запросы по адресам из списка
// подсчитать количество для каждого StatusCode ответа
// предусмотреть возможность отмены запроса по таймауту
func process(urls []string) map[int]int {
	statusCodeCounts := make(map[int]int, len(urls))
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	for _, url := range urls {
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				fmt.Println("Some error ", err.Error())
			}

			response, err := client.Do(req)
			if err != nil {
				fmt.Println("Error in response ", err.Error())
				return
			}
			mu.Lock()
			defer mu.Unlock()
			statusCodeCounts[response.StatusCode]++
		}()
	}

	wg.Wait()
	return statusCodeCounts
}
