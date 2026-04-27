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

// TestStoreObject is responsible for testing
// the storeObject() function.
func TestStoreObject(t *testing.T) {
	a := assert.New(t)

	// Creating the database and the directory
	tempDir := t.TempDir()
	db, err := NewDatabase(tempDir)
	a.Nil(err)
	a.NotNil(db)

	// Creating a Manifest object.
	blocks := []string{"vmnerovneq", "kcqervqerv", "qvqnervnerveqrv"}
	mani := newManifest("filename", "filedir", int64(34), blocks)
	a.NotNil(mani)

	// Creating an IndexTable object.
	it := newIndexTable()
	a.NotNil(it)

	// Storing the indexTable in the database.
	err = db.storeObject("indexTable", it)
	a.Nil(err)

	// Storing the Manifest in the database.
	err = db.storeObject(mani.FileName, mani)
	a.Nil(err)

}