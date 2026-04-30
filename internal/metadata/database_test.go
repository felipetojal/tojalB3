package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewDatabase is responsible for testing the 
// NewDatabase() function.
func TestNewDatabase(t *testing.T) {
	assert := assert.New(t)

	// Creates a temporary directory in-memory
	// for the test.
	tempDir := t.TempDir()
	
	// Calling the function to be tested and
	// verifying its result.
	db, err := NewDatabase(tempDir)
	assert.NotNil(db)
	assert.Nil(err)
	defer db.db.Close()
}

