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

// TestStoreManfiest is responsible for testing
// the StoreManifest function.
func TestStoreManifest(t *testing.T) {
	a := assert.New(t)

	// Creating the temporary directory.
	tempDir := t.TempDir()

	// Creating the database.
	db, err := NewDatabase(tempDir)
	a.NotNil(db)
	a.Nil(err)

	// Creating the manifest to be stored.
	fileName := "opa.png"
	fileDir := "/josias/"
	blocks := []string{"mvwe", "pemw", "egvwoerng-", "gmkwenrg2455g24"}
	fileSize := int64(34)
	m := newManifest(fileName, fileDir, fileSize, blocks)
	a.NotNil(m)

	// Storing the Manifest.
	err = db.StoreManifest(m)
	a.Nil(err)
}

func TestLoadManifest(t *testing.T) {
	a := assert.New(t)

	// Creating the temporary directory.
	tempDir := t.TempDir()

	// Creating the database.
	db, err := NewDatabase(tempDir)
	a.NotNil(db)
	a.Nil(err)

	// Creating the manifest to be stored.
	fileName := "opa.png"
	fileDir := "/josias/"
	blocks := []string{"mvwe", "pemw", "egvwoerng-", "gmkwenrg2455g24"}
	fileSize := int64(34)
	m := newManifest(fileName, fileDir, fileSize, blocks)
	a.NotNil(m)

	// Storing the Manifest.
	err = db.StoreManifest(m)
	a.Nil(err)

	mani, err := db.LoadManifest(fileName)
	a.NotNil(mani)
	a.Nil(err)
}

func TestStoreIndexTable(t *testing.T) {
	a := assert.New(t)

	// Creating the temporary directory.
	tempDir := t.TempDir()

	// Creating the database.
	db, err := NewDatabase(tempDir)
	a.NotNil(db)
	a.Nil(err)

	// Creating the index table
	it := newIndexTable()
	a.NotNil(it)

	// Storing the index table
	err = db.StoreIndexTable(it)
	a.Nil(err)
}

func TestLoadIndexTable(t *testing.T) {
	a := assert.New(t)

	// Creating the temporary directory.
	tempDir := t.TempDir()

	// Creating the database.
	db, err := NewDatabase(tempDir)
	a.NotNil(db)
	a.Nil(err)

	// Creating the index table
	it := newIndexTable()
	a.NotNil(it)

	// Storing the index table
	err = db.StoreIndexTable(it)
	a.Nil(err)

	// Retrieving the index table from the database.
	it2, err := db.LoadIndexTable()
	a.NotNil(it2)
	a.Nil(err)
}
