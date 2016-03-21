package utils

import "crypto/sha256"
import "io"
import "fmt"

// Hash ...
func Hash(password string) string {
	h := sha256.New()
	io.WriteString(h, password)
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

// Compare ...
func Compare(password string, hashedPassword string) bool {
	s := Hash(password)
	if s == hashedPassword {
		return true
	}
	return false
}