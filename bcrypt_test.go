package main

import (
    "fmt"
    "testing"
    "golang.org/x/crypto/bcrypt"
)

func TestBcryptCompare(t *testing.T) {
    hashedPassword := "$2a$10$u9wc1IxoF6uGo2aQYvVTQu20LjJnCLj8Aa3si3AxBANLcb0VKKa9W" // Hash dari database
    password := "test123" // Password asli

    fmt.Printf("Original password: %q\n", password)
    fmt.Printf("Database hash: %q\n", hashedPassword)

    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        t.Errorf("Password comparison failed: %v", err)
    } else {
        t.Log("Password comparison succeeded")
    }

    // Generate new hash for debugging
    newHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Printf("Newly generated hash: %s\n", string(newHash))
}
