package engine

import (
	"crypto/sha256"
	"fmt"
)

// This function takes in the bytes in the block,
// calculates the hash and then return it.
func generateHash(block []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(block)
	if err != nil {
		return "", fmt.Errorf("error generating block hash: %w", err)
	}

	hash := h.Sum(nil)

	return string(hash), nil
}
