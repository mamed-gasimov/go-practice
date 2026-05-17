package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// 1. Extend the API to allow for array of nums and to return the sequences
// 2. Compute the sequences concurrently
// 3. Implement a way to stop as soon as the first routine returns and return the winner
// 4. Store the results of every calculation in a map
// 5. Use the map for concurrent lookup of the solutions
// 6. Handle user interruption
type req struct {
	Num int
}
type resp struct {
	Sequence []int
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/collatz", func(w http.ResponseWriter, r *http.Request) {
		in, _ := io.ReadAll(r.Body)
		data := &req{}
		err := json.Unmarshal(in, data)
		if err != nil {
			log.Fatalf(err.Error())
		}
		res := &resp{Sequence: collatz(data.Num)}
		out, err := json.Marshal(res)
		w.Write(out)
	})
	http.ListenAndServe(":8080", mux)
}
func collatz(n int) []int {
	iter := []int{}
	for n != 1 {
		iter = append(iter, n)
		if n%2 == 0 {
			n = n / 2
			continue
		}
		n = 3*n + 1
	}
	return iter
}
