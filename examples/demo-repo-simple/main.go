package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// Database simulates a database connection
type Database struct {
	users map[int]User
	mu    sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[int]User),
	}
}

func (db *Database) AddUser(user User) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.users[user.ID] = user
}

func (db *Database) GetUser(id int) (User, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	user, exists := db.users[id]
	return user, exists
}

// String processing functions
func processStrings(data []string) []string {
	var result []string
	for _, s := range data {
		// Use strings.Builder for efficient string concatenation
		var sb strings.Builder
		sb.WriteString("processed_")
		sb.WriteString(s)
		sb.WriteString("_end")
		result = append(result, sb.String())
	}
	return result
}

func generateRandomData(size int) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(rand.Intn(256))
	}
	return data
}

func processJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func isValidEmail(email string) bool {
	// Simple email validation using regex
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func main() {
	// Initialize database
	db := NewDatabase()

	// Add some initial data
	db.AddUser(User{ID: 1, Username: "alice", Email: "alice@example.com", CreatedAt: time.Now().Format(time.RFC3339)})
	db.AddUser(User{ID: 2, Username: "bob", Email: "bob@example.com", CreatedAt: time.Now().Format(time.RFC3339)})

	fmt.Println("Demo application started")
}
