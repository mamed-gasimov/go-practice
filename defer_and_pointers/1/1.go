package main

import "fmt"

func main() {
	x := 10

	defer fmt.Println("A:", x)
	defer func() { fmt.Println("B:", x) }()

	x = 42
}

// defer с аргументом захватывает значение сразу (x=10).
// Замыкание читает x в момент выполнения (x=42).
// Defer выполняется LIFO.
// B: 42
// A: 10
