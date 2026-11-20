package utils

import (
	"fmt"
	"os"
)

// Getenv returns the value of key, or fallback when unset or empty.
func Getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// MustGetenv returns the value of key or panics when unset or empty.
func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return v
}
