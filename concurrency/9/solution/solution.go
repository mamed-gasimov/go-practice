package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// имеется функция, которая работает неопределенно долго (до 100 секунд)
func randomTimeWork() {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
}

// написать обертку для этой функции, которая будет прерывать выполнение, если
// функция работает дольше 3 секунд, и возвращать ошибку
func predictableTimeWork() error {
	ch := make(chan struct{})

	go func() {
		randomTimeWork()
		close(ch)
	}()

	select {
	case <-time.After(3 * time.Second):
		return errors.New("Exits after 3 seconds")
	case <-ch:
		fmt.Println("Exits in less than 3 seconds")
		return nil
	}
}

func main() {
	err := predictableTimeWork()
	if err != nil {
		fmt.Println(err.Error())
	}
}
