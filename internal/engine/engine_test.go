package engine

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipetojal/tojalB3/internal/metadata"
	"github.com/felipetojal/tojalB3/internal/volume"
)

// Auxiliary function to create a new test engine.
func setupTestEngine(t *testing.T) (*Engine, string) {
	tempDir := t.TempDir()
	t.Helper()

	// Creating the database.
	dbPath := filepath.Join(tempDir, "badger_db")
	db, err := metadata.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	// Creating the volume manager.
	volPath := filepath.Join(tempDir, "volume.dat")
	vol := volume.NewVolumeManager(volPath)

	// Loading the index table from the database.
	it, err := db.LoadIndexTable()
	if err != nil {
		t.Fatalf("failed to load index table: %v", err)
	}

	// Creating the engine.
	newEngine, err := NewEngine(vol, db, it)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	return newEngine, tempDir
}
func TestStoreFile(t *testing.T) {
	eng, tempDir := setupTestEngine(t)

	// Creating some content (bytes) to store.
	content := bytes.Repeat([]byte("A"), 5000)

	// Creating file1
	file1Path := filepath.Join(tempDir, "file1.txt")
	os.WriteFile(file1Path, content, 0644)

	// Creating file2 (different name, but IDENTICAL content)
	file2Path := filepath.Join(tempDir, "file2.txt")
	os.WriteFile(file2Path, content, 0644)

	// Uploading file1 to the engine.
	err := eng.StoreFile(file1Path)
	if err != nil {
		t.Fatalf("failed storing file1: %v", err)
	}

	// Check the state after the first file
	if len(eng.it.Indexes) != 2 {
		t.Errorf("expected 2 blocks in IndexTable, found %d", len(eng.it.Indexes))
	}

	// Get the hash of the first block to check the RefCount
	var firstHash string
	for hash := range eng.it.Indexes {
		firstHash = hash
		break
	}

	if eng.it.Indexes[firstHash].RefCount != 1 {
		t.Errorf("expected RefCount 1, found %d", eng.it.Indexes[firstHash].RefCount)
	}

	// Upload the SECOND file (Deduplication comes into play)
	err = eng.StoreFile(file2Path)
	if err != nil {
		t.Fatalf("unexpected failure storing file 2: %v", err)
	}

	// Deduplication check.
	// The number of blocks should not have grown, it should still be 2.
	if len(eng.it.Indexes) != 2 {
		t.Errorf("deduplication failed! Expected to keep 2 blocks, found %d", len(eng.it.Indexes))
	}

	// The RefCount MUST be 2 now.
	if eng.it.Indexes[firstHash].RefCount != 2 {
		t.Errorf("deduplication failed! Expected RefCount 2, found %d", eng.it.Indexes[firstHash].RefCount)
	}

	// Verify if the Manifests were stored correctly in BadgerDB
	m1, _ := eng.d.LoadManifest(file1Path)
	m2, _ := eng.d.LoadManifest(file2Path)

	if m1 == nil || m2 == nil {
		t.Errorf("the manifests were not saved in the database")
	}

	if len(m1.Blocks) != 2 || len(m2.Blocks) != 2 {
		t.Errorf("the manifests should have 2 hashes in the list")
	}
}

func TestReadBlock(t *testing.T) {
	// Create a temporary directory for the test files
	tempDir := t.TempDir()

	// Define our test scenarios
	tests := []struct {
		name           string
		contentSize    int
		expectedEOF    bool
		expectedBuffer []byte // How we expect the buffer to look at the end
	}{
		{
			name:        "Scenario 1: File smaller than the block (Padding required)",
			contentSize: 10,
			expectedEOF: true,
		},
		{
			name:        "Scenario 2: File the exact size of the block",
			contentSize: block_size,
			expectedEOF: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE: Prepare the file with the specific size for this scenario
			filePath := filepath.Join(tempDir, "test_read.txt")
			content := bytes.Repeat([]byte("A"), tt.contentSize)
			os.WriteFile(filePath, content, 0644)

			f, err := os.Open(filePath)
			if err != nil {
				t.Fatalf("failed to open file: %v", err)
			}
			defer f.Close()

			buf := make([]byte, block_size)
			endOfFile := false

			// ACT: Call our auxiliary function
			err = readBlock(f, buf, &endOfFile)

			// ASSERT: Check if the error wasn't fatal
			if err != nil && !errors.Is(err, io.EOF) {
				t.Fatalf("unexpected error in readBlock: %v", err)
			}

			// ASSERT: Check if the EOF variable was changed correctly
			if endOfFile != tt.expectedEOF {
				t.Errorf("expected endOfFile = %v, found %v", tt.expectedEOF, endOfFile)
			}

			// ASSERT: Check the Padding (Zero-filling)
			// If the content read is smaller than the block, the rest MUST be byte 0.
			if tt.contentSize < block_size {
				// Check the first padding byte
				if buf[tt.contentSize] != 0 {
					t.Errorf("padding failed! the byte at position %d is not zero", tt.contentSize)
				}
				// Check the last byte of the buffer to ensure everything was cleared
				if buf[block_size-1] != 0 {
					t.Errorf("padding failed! the last byte of the buffer is not zero")
				}
			}
		})
	}
}
