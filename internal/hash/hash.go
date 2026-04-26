package hash

import (
	"crypto/sha256"
	"fmt"
)

// GenerateHash receives the block bytes and return the hash
// in string type if no errors occur.
func GenerateHash(block []byte) (string, error) {
	h := sha256.New()
	n, err := h.Write(block)
	if err != nil {
		return "", fmt.Errorf("error writing hash: %w", err)
	}

	if n == 0 {
		return "", fmt.Errorf("error: 0 bytes were written in hash.")
	}

	hash := h.Sum(nil)

	return fmt.Sprint(string(hash)), nil
}
