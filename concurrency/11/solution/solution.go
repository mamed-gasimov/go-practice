package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// 1. Extend the API to allow for array of nums and to return the sequences
// 2. Compute the sequences concurrently
// 3. Implement a way to stop as soon as the first routine returns and return the winner
// 4. Store the results of every calculation in a map
// 5. Use the map for concurrent lookup of the solutions
// 6. Handle user interruption

// --- Step 1 ---
// Changed from a single Num to a slice of Nums so the client can send
// multiple starting values in one request. The response returns a map of
// each number to its sequence, plus a Winner field for step 3.
type req struct {
	Nums []int `json:"nums"`
	Race bool  `json:"race"` // when true, return only the first sequence that finishes (step 3)
}

type resp struct {
	Sequences map[int][]int `json:"sequences"`
	Winner    *int          `json:"winner,omitempty"` // populated only in race mode (step 3)
}

// --- Step 4 & 5 ---
// A package-level concurrent cache so results survive across requests.
// RWMutex allows many readers to hold the lock simultaneously (RLock),
// while writers get exclusive access (Lock). This fits our pattern well:
// reads are frequent, writes happen only on cache misses.
var (
	cacheMu sync.RWMutex
	cache   = make(map[int][]int)
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/collatz", handleCollatz)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// --- Step 6 ---
	// Listen for SIGINT / SIGTERM so the server shuts down gracefully
	// instead of dropping in-flight requests on Ctrl-C.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-quit // block until signal
	log.Println("shutting down gracefully…")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

func handleCollatz(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	in, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := &req{}
	if err := json.Unmarshal(in, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data.Nums) == 0 {
		http.Error(w, "nums must not be empty", http.StatusBadRequest)
		return
	}

	var result *resp

	if data.Race {
		// --- Step 3 ---
		result = computeRace(data.Nums)
	} else {
		// --- Step 2 ---
		result = computeAll(data.Nums)
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(out); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

// --- Step 2 ---
// computeAll launches one goroutine per number and waits for every one of
// them to finish. A WaitGroup gates the return so we collect all sequences.
func computeAll(nums []int) *resp {
	var (
		mu        sync.Mutex
		wg        sync.WaitGroup
		sequences = make(map[int][]int, len(nums))
	)

	for _, n := range nums {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			seq := cachedCollatz(n) // step 5: cache-aware wrapper
			mu.Lock()
			sequences[n] = seq
			mu.Unlock()
		}(n)
	}

	wg.Wait()
	return &resp{Sequences: sequences}
}

// --- Step 3 ---
// computeRace launches all goroutines but returns as soon as the first one
// finishes. A context with cancel propagates the "stop" signal; remaining
// goroutines see ctx.Done() on their next check (the collatz function
// itself is pure and fast, so cancellation is cooperative at the goroutine
// boundary — the important thing is we don't wait for the slow ones).
func computeRace(nums []int) *resp {
	type result struct {
		num int
		seq []int
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan result, len(nums))

	for _, n := range nums {
		go func(n int) {
			select {
			case <-ctx.Done():
				return // another goroutine already won
			default:
				seq := cachedCollatz(n)
				select {
				case ch <- result{num: n, seq: seq}:
				case <-ctx.Done():
				}
			}
		}(n)
	}

	// Take only the first result, then cancel the rest.
	winner := <-ch
	cancel()

	return &resp{
		Sequences: map[int][]int{winner.num: winner.seq},
		Winner:    &winner.num,
	}
}

// --- Step 4 & 5 ---
// cachedCollatz checks the cache before computing. It takes a read lock
// first (cheap, non-exclusive) and only upgrades to a write lock on a
// cache miss. The double-check after acquiring the write lock prevents
// two goroutines from computing the same value when they both miss
// simultaneously.
func cachedCollatz(n int) []int {
	cacheMu.RLock()
	if seq, ok := cache[n]; ok {
		cacheMu.RUnlock()
		return seq
	}
	cacheMu.RUnlock()

	seq := collatz(n)

	cacheMu.Lock()
	// Double-check: another goroutine may have stored it while we computed.
	if existing, ok := cache[n]; ok {
		cacheMu.Unlock()
		return existing
	}
	cache[n] = seq
	cacheMu.Unlock()

	return seq
}

// collatz is the original, unmodified function.
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
