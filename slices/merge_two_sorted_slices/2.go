package main

import "fmt"

func mergeTwoSortedSlices(a []int, b []int) []int {
	result := make([]int, 0, len(a)+len(b))

	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}

	result = append(result, a[i:]...)
	result = append(result, b[j:]...)

	return result
}

func main() {
	a := []int{1, 2, 3, 4, 5}
	b := []int{4, 5, 6, 7, 8}
	fmt.Println(mergeTwoSortedSlices(a, b))
}
