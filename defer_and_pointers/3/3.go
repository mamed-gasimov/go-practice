package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		defer fmt.Println("val:", i)
	}
}

// Аргумент i вычисляется при регистрации defer (не замыкание).
// Три defer с i=0,1,2 — выполняются в обратном порядке (LIFO).
// val: 2
// val: 1
// val: 0
