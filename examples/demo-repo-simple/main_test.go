package main

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkProcessStrings(b *testing.B) {
	input := []string{"hello", "world", "test", "benchmark", "performance"}

	for i := 0; i < b.N; i++ {
		processStrings(input)
	}
}

func BenchmarkGenerateRandomData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateRandomData(1024)
	}
}

func BenchmarkProcessJSON(b *testing.B) {
	jsonData := []byte(`{"name": "test", "value": 123, "nested": {"field": "value"}}`)

	for i := 0; i < b.N; i++ {
		processJSON(jsonData)
	}
}

func BenchmarkDatabaseOperations(b *testing.B) {
	db := NewDatabase()

	for i := 0; i < b.N; i++ {
		user := User{ID: i, Username: fmt.Sprintf("user%d", i), Email: fmt.Sprintf("user%d@example.com", i)}
		db.AddUser(user)
		_, _ = db.GetUser(i)
	}
}

func BenchmarkStringProcessing(b *testing.B) {
	data := make([]string, 100)
	for i := range data {
		data[i] = fmt.Sprintf("test_string_%d", i)
	}

	for i := 0; i < b.N; i++ {
		processStrings(data)
	}
}

func BenchmarkJSONProcessing(b *testing.B) {
	jsonData := []byte(`{"id": 1, "name": "test", "email": "test@example.com", "created_at": "` + time.Now().Format(time.RFC3339) + `"}`)

	for i := 0; i < b.N; i++ {
		processJSON(jsonData)
	}
}

func BenchmarkRegexValidation(b *testing.B) {
	emails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"invalid-email",
		"another.test@sub.domain.org",
	}

	for i := 0; i < b.N; i++ {
		for _, email := range emails {
			isValidEmail(email)
		}
	}
}

func BenchmarkMutexOperations(b *testing.B) {
	db := NewDatabase()

	for i := 0; i < b.N; i++ {
		user := User{ID: i % 100, Username: fmt.Sprintf("user%d", i), Email: fmt.Sprintf("user%d@example.com", i)}
		db.AddUser(user)
	}
}
