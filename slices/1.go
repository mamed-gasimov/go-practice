package main

import "fmt"

func appendLen(numbers []*int) {
	size := len(numbers)
	numbers = append(numbers, &size)
	// The local numbers now has length 4. But the caller's numbers
	// still has length 3 — the append modified only the local slice header's length field.
	// However, because the original slice had capacity 5,
	// the append didn't allocate a new backing array.
	// It wrote &size into index 3 of the same underlying array.
	// The caller just can't see it through its slice header (length is still 3).
}

func main() {
	numbers := make([]*int, 0, 5)
	var number int

	// There's only one number variable. Each iteration mutates it and appends its address.
	// After the loop, numbers has 3 elements, but all three pointers point to the same variable number,
	// whose current value is 3. So dereferencing any of them gives 3.
	for range 3 {
		number++
		numbers = append(numbers, &number)
	}
	appendLen(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number) // 3 3 3
	}

	fmt.Println("\n------------------------")

	for _, number := range numbers[:4] {
		fmt.Printf("%d ", *number) // 3 3 3 3
	}

	fmt.Println("\n------------------------")
	double(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number) // 6 6 6
	}
}

func double(numbers []*int) {
	// Since all three slice elements hold the same pointer (&number from main),
	// naively iterating the slice would double that single variable three times (3→6→12→24).
	// The fix is to collect unique pointers first, then double each one exactly once.
	m := map[*int]struct{}{}

	for _, v := range numbers {
		m[v] = struct{}{}
	}

	for k := range m {
		*k = *k * 2
	}
}
