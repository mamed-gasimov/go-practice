package main

import "fmt"

type Counter struct{ n int }

func (c Counter) PrintVal()  { fmt.Println("val:", c.n) }
func (c *Counter) PrintPtr() { fmt.Println("ptr:", c.n) }

func main() {
	c := Counter{n: 1}

	defer c.PrintVal() // val: 1
	defer c.PrintPtr() // ptr: 99

	c.n = 99
}

// defer c.PrintVal() — receiver-значение копируется сразу (c.n=1).
// defer c.PrintPtr() — receiver-указатель, читает c.n в момент выполнения (c.n=99). LIFO.
// ptr: 99
// val: 1
