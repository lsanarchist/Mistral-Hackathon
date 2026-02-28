package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

func main() {
	// Register pprof handlers
	http.HandleFunc("/debug/pprof/", pprof.Index)

	// Performance endpoints
	http.HandleFunc("/cpu-hotspot", cpuHotspotHandler)
	http.HandleFunc("/alloc-heavy", allocHeavyHandler)
	http.HandleFunc("/mutex-contention", mutexContentionHandler)

	fmt.Println("Demo server running on :6060")
	fmt.Println("Endpoints:")
	fmt.Println("- /cpu-hotspot - CPU intensive endpoint")
	fmt.Println("- /alloc-heavy - Memory allocation heavy endpoint")
	fmt.Println("- /mutex-contention - Mutex contention endpoint")
	fmt.Println("- /debug/pprof/ - pprof endpoints")

	log.Fatal(http.ListenAndServe(":6060", nil))
}

func cpuHotspotHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate CPU hotspot
	start := time.Now()
	for i := 0; i < 100000000; i++ {
		_ = i * i
	}
	fmt.Fprintf(w, "CPU hotspot completed in %v\n", time.Since(start))
}

func allocHeavyHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate memory allocation
	for i := 0; i < 10000; i++ {
		buf := make([]byte, 1024*1024) // 1MB per allocation
		_ = buf
	}
	fmt.Fprintf(w, "Allocation heavy completed\n")
}

func mutexContentionHandler(w http.ResponseWriter, r *http.Request) {
	var mu sync.Mutex
	var counter int

	// Create contention
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Fprintf(w, "Mutex contention completed, counter: %d\n", counter)
}