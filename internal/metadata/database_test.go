package metadata

import (
	"path/filepath"
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
	blocks := []string{"mvwe", "pemw", "egvwoerng-", "gmkwenrg2455g24"}
	fileSize := int64(34)
	m := NewManifest(fileName, fileSize)
	m.AddBlock(blocks...)
	a.NotNil(m)

	// Storing the Manifest.
	err = db.StoreManifest(m)
	a.Nil(err)
}

func TestDeleteManifest(t *testing.T) {
	// ARRANGE: Setup temporary database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_badger_db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Prepare a fake manifest
	fileName := "test_delete.txt"
	mani := NewManifest(fileName, 2048)
	mani.AddBlock("hash_123")
	mani.AddBlock("hash_456")

	// Store the manifest first to ensure it exists in the DB
	err = db.StoreManifest(mani)
	if err != nil {
		t.Fatalf("failed to store manifest during setup: %v", err)
	}

	// ACT 1: Delete the EXISTING manifest
	err = db.DeleteManifest(fileName)
	if err != nil {
		t.Errorf("unexpected error deleting existing manifest: %v", err)
	}

	// ASSERT 1: Verify it was actually deleted from BadgerDB
	deletedMani, err := db.LoadManifest(fileName)
	if deletedMani != nil {
		t.Errorf("expected manifest to be nil after deletion, but got data")
	}
	if err == nil {
		t.Errorf("expected an error (like key not found) when loading a deleted manifest, got nil")
	}

	// ACT 2: Delete a NON-EXISTING manifest (Idempotency Test)
	// Since we chose Option A, deleting a ghost file should succeed silently.
	err = db.DeleteManifest("ghost_file.txt")
	if err != nil {
		t.Errorf("expected no error when deleting a non-existing manifest (idempotency), but got: %v", err)
	}
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
	blocks := []string{"mvwe", "pemw", "egvwoerng-", "gmkwenrg2455g24"}
	fileSize := int64(34)
	m := NewManifest(fileName, fileSize)
	m.AddBlock(blocks...)
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
