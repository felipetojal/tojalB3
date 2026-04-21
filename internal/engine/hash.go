package engine

import (
	"crypto/sha256"
	"fmt"
)

// This function takes in bytes, calculates
// the hash and then return it as a string.
func generateHash(block []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(block)
	if err != nil {
		return "", fmt.Errorf("error generating hash: %w", err)
	}

	hash := h.Sum(nil)

	return string(hash), nil
}
