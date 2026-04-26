package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewIndex is responsible for testing
// the Index constructor.
func TestNewIndex(t *testing.T) {
	a := assert.New(t)

	hash := "ansognqeorngqenrgafçvqgpfq"
	address := 34
	i := newIndex(hash, address)
	a.NotNil(i)
}

// TestNewTable is responsible for testing
// the IndexTable constructor.
func TestNewIndexTable(t *testing.T) {
	a := assert.New(t)

	it := newIndexTable()
	a.NotNil(it)
}