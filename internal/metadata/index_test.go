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
	a.Equal(1, i.refCount)
}

// TestNewTable is responsible for testing
// the IndexTable constructor.
func TestNewIndexTable(t *testing.T) {
	a := assert.New(t)

	it := newIndexTable()
	a.NotNil(it)
}

// TestGetIndex is responsible for testing
// the getIndex function associated to IndexTable.
func TestGetIndex(t *testing.T) {
	a := assert.New(t)

	// Creating the Index object.
	hash := "vnwerogvweg23ujg234"
	address := 45
	i := newIndex(hash, address)
	a.NotNil(i)

	// Creating the table and adding the Index.
	it := newIndexTable()
	a.NotNil(it)
	it.addIndex(*i)

	// Getting the Index object.
	i, err := it.getIndex(hash)
	a.NotNil(i)
	a.Nil(err)
}

// TestAddIndex is responsible for testing
// the addIndex function associated to IndexTable.
func TestAddIndex(t *testing.T) {
	a := assert.New(t)

	// Creating the Index object.
	hash := "vnwerogvweg23ujg234"
	address := 45
	i := newIndex(hash, address)
	a.NotNil(i)
	a.Equal(1, i.refCount)

	// Creating the table and adding the Index.
	it := newIndexTable()
	a.NotNil(it)
	it.addIndex(*i)
	savedIndex1 := it.indexes[i.hash]
	a.Equal(1, savedIndex1.refCount)

	// Checking if the hash we just inserted
	// is really there.
	b := checkExistence(it, hash)
	a.True(b)

	// Creating a new index with the same hash.
	address = 9
	i2 := newIndex(hash, address)
	a.NotNil(i2)
	it.addIndex(*i2)
	savedIndex2 := it.indexes[i2.hash]
	a.Equal(2, savedIndex2.refCount)
}

