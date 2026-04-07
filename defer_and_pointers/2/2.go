package main

import "fmt"

func calculate() (result int) {
	defer func() { result *= 2 }()
	result = 5
	return result + 10
	// return result+10 записывает 15 в именованную переменную result.
	// Затем defer умножает result на 2 → 30.
}

func main() {
	fmt.Println(calculate()) // 30
}
