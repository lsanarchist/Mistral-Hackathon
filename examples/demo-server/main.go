package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

// User represents a user entity for realistic data
type User struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Friends   []string `json:"friends"`
	CreatedAt string   `json:"createdAt"`
	Metadata  struct {
		Preferences map[string]interface{} `json:"preferences"`
		Tags       []string               `json:"tags"`
	} `json:"metadata"`
}

// Database simulates a slow database connection
var database = struct {
	mu      sync.Mutex
	users   []User
	queries int
}{
	users: make([]User, 0),
}

func main() {
	// Seed random for realistic data
	rand.Seed(time.Now().UnixNano())
	
	// Initialize demo data
	initializeDemoData()

	// Register pprof handlers
	http.HandleFunc("/debug/pprof/", pprof.Index)

	// Performance endpoints with realistic issues
	http.HandleFunc("/api/users", usersHandler)
	http.HandleFunc("/api/search", searchHandler)
	http.HandleFunc("/api/analytics", analyticsHandler)
	http.HandleFunc("/api/export", exportHandler)
	http.HandleFunc("/api/process", processHandler)
	http.HandleFunc("/api/strings", inefficientStringHandler)
	http.HandleFunc("/api/nocache", noCacheHandler)
	http.HandleFunc("/api/iobound", ioBoundHandler)
	http.HandleFunc("/api/leak", memoryLeakHandler)
	http.HandleFunc("/api/blocking", blockingIOHandler)
	http.HandleFunc("/api/goroutine", goroutineLeakHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("🚀 Enhanced Demo Server running on :6060")
	fmt.Println("📊 Endpoints:")
	fmt.Println("  - GET  /api/users     - User listing with JSON serialization overhead")
	fmt.Println("  - GET  /api/search    - Search with database contention")
	fmt.Println("  - GET  /api/analytics - CPU-intensive analytics processing")
	fmt.Println("  - GET  /api/export    - Memory-heavy data export")
	fmt.Println("  - POST /api/process   - Complex business logic with mutex contention")
	fmt.Println("  - GET  /api/strings   - Inefficient string concatenation")
	fmt.Println("  - GET  /api/nocache   - Lack of caching for expensive computations")
	fmt.Println("  - GET  /api/iobound   - I/O bottleneck simulation")
	fmt.Println("  - GET  /api/leak      - Memory leak simulation")
	fmt.Println("  - GET  /api/blocking  - Blocking I/O operations")
	fmt.Println("  - GET  /api/goroutine - Goroutine leak simulation")
	fmt.Println("  - GET  /health       - Health check endpoint")
	fmt.Println("  - GET  /debug/pprof/  - Performance profiling endpoints")
	fmt.Println("\n🎯 Performance Issues Demonstrated:")
	fmt.Println("  🔥 JSON serialization overhead (/api/users)")
	fmt.Println("  🔥 Database lock contention (/api/search)")
	fmt.Println("  🔥 CPU-bound analytics (/api/analytics)")
	fmt.Println("  🔥 Memory allocation patterns (/api/export)")
	fmt.Println("  🔥 Mutex contention in business logic (/api/process)")
	fmt.Println("  🔥 I/O bottlenecks (/api/iobound)")
	fmt.Println("  🔥 Inefficient string operations (/api/strings)")
	fmt.Println("  🔥 Lack of result caching (/api/nocache)")
	fmt.Println("  🔥 Memory leaks (/api/leak) - ~400MB leaked!")
	fmt.Println("  🔥 Blocking I/O operations (/api/blocking)")
	fmt.Println("  🔥 Goroutine leaks (/api/goroutine) - 200 leaking goroutines!")

	log.Fatal(http.ListenAndServe(":6060", nil))
}

func initializeDemoData() {
	// Generate realistic user data
	for i := 0; i < 1000; i++ {
		user := User{
			ID:    i,
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Friends: generateFriends(),
			CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(365)).Format(time.RFC3339),
		}
		user.Metadata.Preferences = map[string]interface{}{
			"theme":        randomTheme(),
			"notifications": randomBool(),
			"language":     randomLanguage(),
		}
		user.Metadata.Tags = generateTags()
		database.users = append(database.users, user)
	}
}

func generateFriends() []string {
	friends := []string{}
	count := rand.Intn(50) + 5
	for i := 0; i < count; i++ {
		friends = append(friends, fmt.Sprintf("friend%d@example.com", rand.Intn(1000)))
	}
	return friends
}

func randomTheme() string {
	themes := []string{"light", "dark", "system", "solarized", "dracula"}
	return themes[rand.Intn(len(themes))]
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomLanguage() string {
	langs := []string{"en", "es", "fr", "de", "ja", "zh"}
	return langs[rand.Intn(len(langs))]
}

func generateTags() []string {
	tags := []string{}
	count := rand.Intn(10) + 2
	allTags := []string{"premium", "active", "inactive", "vip", "new", "beta", "early-adopter"}
	for i := 0; i < count; i++ {
		tags = append(tags, allTags[rand.Intn(len(allTags))])
	}
	return tags
}

// usersHandler demonstrates JSON serialization overhead
func usersHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate database query with lock contention
	database.mu.Lock()
	database.queries++
	users := make([]User, len(database.users))
	copy(users, database.users)
	database.mu.Unlock()
	
	// Add artificial delay to simulate real database query
	time.Sleep(10 * time.Millisecond)
	
	// JSON serialization - this is a common bottleneck
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    users,
		"count":   len(users),
		"timestamp": time.Now().Format(time.RFC3339),
	})
	
	fmt.Printf("📊 Users endpoint served %d users in %v\n", len(users), time.Since(start))
}

// searchHandler demonstrates database contention
func searchHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	query := r.URL.Query().Get("q")
	if query == "" {
		query = "test"
	}
	
	// Simulate expensive search with lock contention
	database.mu.Lock()
	defer database.mu.Unlock()
	
	var results []User
	for _, user := range database.users {
		if containsString(user.Name, query) || containsString(user.Email, query) {
			results = append(results, user)
		}
	}
	
	// Simulate complex search processing
	time.Sleep(25 * time.Millisecond)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"query":   query,
		"results": results,
		"count":   len(results),
	})
	
	fmt.Printf("🔍 Search for '%s' found %d results in %v\n", query, len(results), time.Since(start))
}

func containsString(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 && 
		   (haystack == needle || 
		   (len(haystack) >= len(needle) && haystack[:len(needle)] == needle))
}

// analyticsHandler demonstrates CPU-bound processing
func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate expensive analytics calculation
	var totalFriends int
	var activeUsers int
	var premiumUsers int
	
	for _, user := range database.users {
		totalFriends += len(user.Friends)
		if rand.Float32() > 0.3 {
			activeUsers++
		}
		for _, tag := range user.Metadata.Tags {
			if tag == "premium" {
				premiumUsers++
				break
			}
		}
		
		// Simulate complex calculations
		for i := 0; i < 1000; i++ {
			_ = i * i * rand.Intn(100)
		}
	}
	
	// More CPU-intensive work
	fibonacci(30)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"totalUsers":      len(database.users),
		"activeUsers":     activeUsers,
		"premiumUsers":    premiumUsers,
		"avgFriends":      float64(totalFriends) / float64(len(database.users)),
		"processingTime":  time.Since(start).String(),
	})
	
	fmt.Printf("📈 Analytics processed in %v\n", time.Since(start))
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// exportHandler demonstrates memory allocation patterns
func exportHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate memory-heavy export operation
	var exportData []map[string]interface{}
	
	for _, user := range database.users {
		userCopy := make(map[string]interface{})
		userCopy["id"] = user.ID
		userCopy["name"] = user.Name
		userCopy["email"] = user.Email
		
		// Create large metadata copies - this is inefficient!
		friendsCopy := make([]string, len(user.Friends))
		copy(friendsCopy, user.Friends)
		userCopy["friends"] = friendsCopy
		
		prefsCopy := make(map[string]interface{})
		for k, v := range user.Metadata.Preferences {
			prefsCopy[k] = v
		}
		userCopy["preferences"] = prefsCopy
		
		exportData = append(exportData, userCopy)
		
		// Simulate memory pressure - allocate 3MB per user (increased from 2MB)
		largeBuffer := make([]byte, 3*1024*1024)
		_ = largeBuffer
		
		// Add additional memory allocation patterns
		for i := 0; i < 5; i++ {
			tempSlice := make([]int, 100000)  // Additional 800KB allocation
			_ = tempSlice
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"exportSize": len(exportData),
		"data":       exportData,
	})
	
	fmt.Printf("💾 Export generated %d records in %v\n", len(exportData), time.Since(start))
}

// inefficientStringHandler demonstrates poor string concatenation
func inefficientStringHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Inefficient string concatenation - classic performance issue
	var result string
	for i := 0; i < 1000; i++ {
		result += fmt.Sprintf("User %d: %s\n", i, database.users[i%len(database.users)].Name)
	}
	
	// More inefficient string operations
	for _, user := range database.users {
		for _, friend := range user.Friends {
			result += "Friend: " + friend + "\n"
		}
	}
	
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
	
	fmt.Printf("📝 Inefficient string handler completed in %v\n", time.Since(start))
}

// noCacheHandler demonstrates lack of caching
func noCacheHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate expensive computation that should be cached
	var totalOperations int
	
	for _, user := range database.users {
		// Complex calculation that doesn't change
		complexResult := 0
		for i := 0; i < 1000; i++ {
			complexResult += i * len(user.Name) * len(user.Email)
		}
		totalOperations += complexResult
	}
	
	// This same calculation happens on every request - no caching!
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"totalOperations": totalOperations,
		"processingTime": time.Since(start).String(),
	})
	
	fmt.Printf("🔢 No-cache handler processed in %v\n", time.Since(start))
}

// ioBoundHandler demonstrates I/O bottlenecks
func ioBoundHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate I/O bound operations
	var results []string
	
	for i := 0; i < 100; i++ {
		// Simulate file I/O with random delays
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		results = append(results, fmt.Sprintf("I/O operation %d completed", i))
	}
	
	// More I/O simulation
	for range database.users[:50] {
		// Simulate database writes
		time.Sleep(5 * time.Millisecond)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"operations": results,
		"count": len(results),
	})
	
	fmt.Printf("💾 I/O bound handler completed in %v\n", time.Since(start))
}

// processHandler demonstrates mutex contention in business logic
func processHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse request body
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Simulate business logic with mutex contention
	var mu sync.Mutex
	var sharedCounter int
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Simulate work
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			
			// Contended section
			mu.Lock()
			sharedCounter++
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"processed":     sharedCounter,
		"requestData":   requestData,
		"processingTime": time.Since(start).String(),
	})
	
	fmt.Printf("🔄 Process handler completed with %d operations in %v\n", sharedCounter, time.Since(start))
}

// memoryLeakHandler demonstrates memory leak pattern
func memoryLeakHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate memory leak - growing slice that's never released
	var leakySlice [][]byte
	
	for i := 0; i < 200; i++ {  // Increased from 100 to 200 for more pronounced leak
		// Allocate 2MB chunks and keep references (increased from 1MB)
		chunk := make([]byte, 2*1024*1024)
		for j := range chunk {
			chunk[j] = byte(i % 256)
		}
		leakySlice = append(leakySlice, chunk)
		
		// Simulate some processing
		time.Sleep(2 * time.Millisecond)  // Reduced delay for faster leak
	}
	
	// leakySlice is never released - this creates a memory leak
	_ = leakySlice
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Memory leak simulation completed",
		"processingTime": time.Since(start).String(),
	})
	
	fmt.Printf("💀 Memory leak handler completed in %v (leaked ~100MB)\n", time.Since(start))
}

// blockingIOHandler demonstrates blocking I/O operations
func blockingIOHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate blocking I/O operations that tie up goroutines
	for i := 0; i < 100; i++ {  // Increased from 50 to 100 for more blocking
		// Simulate slow I/O with longer random delays
		delay := time.Duration(rand.Intn(100)+20) * time.Millisecond  // Increased delay range
		time.Sleep(delay)
		
		// Simulate more processing after I/O
		for j := 0; j < 2000; j++ {  // Increased computation
			_ = j * j * rand.Intn(10)
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"operations": 50,
		"processingTime": time.Since(start).String(),
	})
	
	fmt.Printf("🚧 Blocking I/O handler completed in %v\n", time.Since(start))
}

// goroutineLeakHandler demonstrates goroutine leak pattern
func goroutineLeakHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Simulate goroutine leak - starting goroutines without waiting
	for i := 0; i < 200; i++ {  // Increased from 100 to 200 for more pronounced leak
		go func(id int) {
			// Simulate work that never completes
			for {
				// Busy wait - this goroutine will never exit
				time.Sleep(50 * time.Millisecond)  // Reduced sleep for more activity
				_ = id * rand.Intn(1000)
				// Add some computation to make it more visible in CPU profiles
				for j := 0; j < 100; j++ {
					_ = j * j * id
				}
			}
		}(i)
		
		// Small delay to make the leak more visible
		time.Sleep(1 * time.Millisecond)  // Reduced delay for faster leak
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Goroutine leak simulation started (100 leaking goroutines)",
		"processingTime": time.Since(start).String(),
	})
	
	fmt.Printf("🌀 Goroutine leak handler completed in %v (leaked 100 goroutines)\n", time.Since(start))
}

// healthHandler provides server health information
func healthHandler(w http.ResponseWriter, r *http.Request) {
	database.mu.Lock()
	queries := database.queries
	database.mu.Unlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(time.Now().Add(-time.Hour)).String(), // Simulated
		"queries":   queries,
		"users":     len(database.users),
	})
}