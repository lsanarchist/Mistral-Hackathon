package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

// Post represents a blog post
type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	AuthorID  int    `json:"author_id"`
	CreatedAt string `json:"created_at"`
}

// Database simulates a database connection
type Database struct {
	users map[int]User
	posts map[int]Post
	mu    sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		users: make(map[int]User),
		posts: make(map[int]Post),
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

func (db *Database) AddPost(post Post) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.posts[post.ID] = post
}

func (db *Database) GetPost(id int) (Post, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	post, exists := db.posts[id]
	return post, exists
}

// APIHandler handles HTTP requests
type APIHandler struct {
	db *Database
}

func NewAPIHandler(db *Database) *APIHandler {
	return &APIHandler{db: db}
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/users":
		h.handleUsers(w, r)
	case "/posts":
		h.handlePosts(w, r)
	default:
		handleNotFound(w, r)
	}
}

func (h *APIHandler) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getUsers(w, r)
	case "POST":
		h.createUser(w, r)
	default:
		handleMethodNotAllowed(w, r)
	}
}

func (h *APIHandler) handlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getPosts(w, r)
	case "POST":
		h.createPost(w, r)
	default:
		handleMethodNotAllowed(w, r)
	}
}

func (h *APIHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	db := h.db
	db.mu.Lock()
	defer db.mu.Unlock()

	var users []User
	for _, user := range db.users {
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *APIHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		handleBadRequest(w, r, err)
		return
	}

	// Validate email format
	if !isValidEmail(user.Email) {
		handleBadRequest(w, r, fmt.Errorf("invalid email format"))
		return
	}

	h.db.AddUser(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *APIHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	db := h.db
	db.mu.Lock()
	defer db.mu.Unlock()

	var posts []Post
	for _, post := range db.posts {
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *APIHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		handleBadRequest(w, r, err)
		return
	}

	// Validate post content
	if len(post.Content) < 10 {
		handleBadRequest(w, r, fmt.Errorf("post content too short"))
		return
	}

	h.db.AddPost(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func isValidEmail(email string) bool {
	// Simple email validation using regex
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found: %s", r.URL.Path)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprintf(w, "Method Not Allowed: %s", r.Method)
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Bad Request: %v", err)
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

func main() {
	// Initialize database
	db := NewDatabase()

	// Add some initial data
	db.AddUser(User{ID: 1, Username: "alice", Email: "alice@example.com", CreatedAt: time.Now().Format(time.RFC3339)})
	db.AddUser(User{ID: 2, Username: "bob", Email: "bob@example.com", CreatedAt: time.Now().Format(time.RFC3339)})

	db.AddPost(Post{ID: 1, Title: "First Post", Content: "This is the first blog post", AuthorID: 1, CreatedAt: time.Now().Format(time.RFC3339)})
	db.AddPost(Post{ID: 2, Title: "Second Post", Content: "This is the second blog post", AuthorID: 2, CreatedAt: time.Now().Format(time.RFC3339)})

	// Set up API handler
	handler := NewAPIHandler(db)

	// Start HTTP server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
