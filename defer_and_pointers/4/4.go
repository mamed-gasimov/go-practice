package main

import "fmt"

func safeDiv(a, b int) (res int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
			res = -1
		}
	}()
	res = a / b
	return
}

func main() {
	fmt.Println(safeDiv(10, 2))
	fmt.Println(safeDiv(10, 0))
}

// Первый вызов — деление без panic, res=5.
// Второй — panic при делении на 0, recover перехватывает, устанавливает res=-1.
// 5
// recovered: runtime error: integer divide by zero
// -1
