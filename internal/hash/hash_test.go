package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateHash is responsible for testing
// GenerateHash function.
func TestGenerateHash(t *testing.T) {
	assert := assert.New(t)

	// Testing the success case.
	b1 := []byte{34, 23, 1, 0, 45, 6, 222, 41, 203, 32, 53, 95, 184, 39}
	s1, err := GenerateHash(b1)
	assert.Nil(err)
	assert.NotEmpty(s1)

	// Testing the error case.
	b2 := []byte{}
	s2, err := GenerateHash(b2)
	assert.NotNil(err)
	assert.Empty(s2)
	
}