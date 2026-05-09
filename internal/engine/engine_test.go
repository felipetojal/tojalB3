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
	"github.com/stretchr/testify/assert"
)

// setupTestEngine is a helper function to create a clean Engine
// inside a disposable temporary directory.
func setupTestEngine(t *testing.T) (*Engine, string) {
	tempDir := t.TempDir()

	// 1. Setup BadgerDB
	dbPath := filepath.Join(tempDir, "badger_db")
	db, err := metadata.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	// 2. Setup Volume
	volPath := filepath.Join(tempDir, "volume.dat")
	vol := volume.NewVolumeManager(volPath)

	// 3. Setup Index Table
	it, err := db.LoadIndexTable()
	if err != nil {
		t.Fatalf("failed to load index table: %v", err)
	}

	// 4. Mount the Engine
	eng, err := NewEngine(vol, db, it)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	return eng, tempDir
}

func TestGetFile(t *testing.T) {
	// Scenario: file does not exist in volume
	t.Run("file does not exist", func(t *testing.T) {
		eng, tempDir := setupTestEngine(t)
		fileName := "file1.txt"
		err := eng.GetFile(fileName, tempDir)
		assert.Error(t, err, "expected error when getting non-existent file %s", fileName)
	})

	// Scenario: test padding (File smaller than 1 block)
	t.Run("Success with Padding Logic", func(t *testing.T) {
		eng, tempDir := setupTestEngine(t)

		// Creating file with content smaller than 1 block
		// and asserting its creation.
		content := []byte("hello world.")
		fileName := "small.txt"
		filePath := filepath.Join(tempDir, fileName)
		f, err := os.Create(filePath)
		assert.NoError(t, err)
		assert.NotNil(t, f, "failed to create small file")
		defer f.Close()
		n, err := f.Write(content)
		assert.NoError(t, err, "failed to write small file")
		assert.Equal(t, len(content), n, "failed to write small file")

		// Storing the file in the databse
		err = eng.StoreFile(filePath)
		assert.NoError(t, err, "failed to store small file")

		// Now it´s time to retrieve the file.
		destDirName := filepath.Join(tempDir, "downloads")
		err = os.Mkdir(destDirName, 0755)
		assert.NoError(t, err, "failed to create downloads directory")
		err = eng.GetFile(filePath, destDirName)
		assert.NoError(t, err, "failed to get file")

		destFilePath := filepath.Join(destDirName, fileName)
		destFile, err := os.Open(destFilePath)
		assert.NoError(t, err, "failed to open destination file")
		assert.NotNil(t, destFile, "failed to open destination file")
		defer destFile.Close()

		// Asserting that the bytes saved were also read.
		bytesRead, err := os.ReadFile(destFilePath)
		assert.NoError(t, err, "failed to read destination file")
		assert.Equal(t, bytesRead, content, "failed to read destination file")
	})

	// Scenario: testing with multiple blocks
	t.Run("Success with Multiple Blocks", func(t *testing.T) {
		eng, tempDir := setupTestEngine(t)

		// Creating file with content smaller than 1 block
		// and asserting its creation.
		content := bytes.Repeat([]byte("hello world."), 10000)
		fileName := "big.txt"
		filePath := filepath.Join(tempDir, fileName)
		f, err := os.Create(filePath)
		assert.NoError(t, err)
		assert.NotNil(t, f, "failed to create big file")
		defer f.Close()
		n, err := f.Write(content)
		assert.NoError(t, err, "failed to write big file")
		assert.Equal(t, len(content), n, "failed to write big file")

		// Storing the file in the databse
		err = eng.StoreFile(filePath)
		assert.NoError(t, err, "failed to store big file")

		// Now it´s time to retrieve the file.
		destDirName := filepath.Join(tempDir, "downloads")
		err = os.Mkdir(destDirName, 0755)
		assert.NoError(t, err, "failed to create downloads directory")
		err = eng.GetFile(filePath, destDirName)
		assert.NoError(t, err, "failed to get file")

		destFilePath := filepath.Join(destDirName, fileName)
		destFile, err := os.Open(destFilePath)
		assert.NoError(t, err, "failed to open destination file")
		assert.NotNil(t, destFile, "failed to open destination file")
		defer destFile.Close()

		// Asserting that the bytes saved were also read.
		bytesRead, err := os.ReadFile(destFilePath)
		assert.NoError(t, err, "failed to read destination file")
		assert.Equal(t, bytesRead, content, "failed to read destination file")
	})
}

func TestListFiles(t *testing.T) {
	eng, tempDir := setupTestEngine(t)

	// Each file that will be stored need to exist and have some content.
	content := []byte("Hello World")
	storedFiles := []string{filepath.Join(tempDir, "file1.txt"), filepath.Join(tempDir, "file2.txt"), filepath.Join(tempDir, "file3.txt")}
	for _, file := range storedFiles {
		// Creating the files and adding content.
		err := os.WriteFile(file, content, 0644)
		assert.NoError(t, err, "failed to create file in disk: %s", file)
		// Now that the files were created, we store them.
		err = eng.StoreFile(file)
		assert.NoError(t, err, "failed to store file %s", file)
	}

	listFiles := eng.ListFiles()
	assert.ElementsMatch(t, listFiles, storedFiles)
}

func TestStoreFile(t *testing.T) {
	eng, tempDir := setupTestEngine(t)

	// Create content (5000 bytes forces the file to be split into 2 blocks)
	content := bytes.Repeat([]byte("A"), 5000)

	file1Path := filepath.Join(tempDir, "file1.txt")
	os.WriteFile(file1Path, content, 0644)

	file2Path := filepath.Join(tempDir, "file2.txt")
	os.WriteFile(file2Path, content, 0644) // Exact same content

	// ACT 1: Upload file1
	err := eng.StoreFile(file1Path)
	if err != nil {
		t.Fatalf("failed storing file1: %v", err)
	}

	// ASSERT 1: Verify the state after the first file
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

	// ACT 2: Upload the SECOND file (Deduplication comes into play)
	err = eng.StoreFile(file2Path)
	if err != nil {
		t.Fatalf("unexpected failure storing file 2: %v", err)
	}

	// ASSERT 2: Deduplication check
	if len(eng.it.Indexes) != 2 {
		t.Errorf("deduplication failed! Expected to keep 2 blocks, found %d", len(eng.it.Indexes))
	}
	if eng.it.Indexes[firstHash].RefCount != 2 {
		t.Errorf("deduplication failed! Expected RefCount 2, found %d", eng.it.Indexes[firstHash].RefCount)
	}

	// ACT 3: Try to upload file1 AGAIN (Should return an error)
	err = eng.StoreFile(file1Path)
	if err == nil {
		t.Errorf("expected an error when trying to store a file that already exists")
	}
}

func TestDeleteFile(t *testing.T) {
	t.Run("File Not Found", func(t *testing.T) {
		eng, _ := setupTestEngine(t)
		err := eng.DeleteFile("ghost_file.txt")
		if err == nil {
			t.Fatal("expected an error when deleting a non-existent file, got nil")
		}
	})

	t.Run("Single File Deletion", func(t *testing.T) {
		eng, tempDir := setupTestEngine(t)
		content := bytes.Repeat([]byte("B"), block_size) // Exactly 1 block
		filePath := filepath.Join(tempDir, "file_to_delete.txt")
		os.WriteFile(filePath, content, 0644)

		eng.StoreFile(filePath)

		// Delete the file
		err := eng.DeleteFile(filePath)
		if err != nil {
			t.Fatalf("failed to delete file: %v", err)
		}

		// Manifest and IndexTable must be clean
		mani, _ := eng.d.LoadManifest(filePath)
		if mani != nil {
			t.Errorf("expected manifest to be deleted from database, but it still exists")
		}
		if len(eng.it.Indexes) != 0 {
			t.Errorf("expected IndexTable to be empty, got %d blocks left", len(eng.it.Indexes))
		}
	})

	t.Run("Deduplication Deletion", func(t *testing.T) {
		eng, tempDir := setupTestEngine(t)
		content := bytes.Repeat([]byte("C"), block_size)

		file1 := filepath.Join(tempDir, "file1.txt")
		file2 := filepath.Join(tempDir, "file2.txt")
		os.WriteFile(file1, content, 0644)
		os.WriteFile(file2, content, 0644)

		eng.StoreFile(file1)
		eng.StoreFile(file2)

		var sharedHash string
		for hash := range eng.it.Indexes {
			sharedHash = hash
			break
		}

		// Delete only the first file
		err := eng.DeleteFile(file1)
		if err != nil {
			t.Fatalf("failed to delete file1: %v", err)
		}

		// RefCount should drop to 1, but block stays
		if eng.it.Indexes[sharedHash].RefCount != 1 {
			t.Errorf("expected RefCount to drop to 1, got %d", eng.it.Indexes[sharedHash].RefCount)
		}

		m1, _ := eng.d.LoadManifest(file1)
		if m1 != nil {
			t.Errorf("expected file1 manifest to be deleted")
		}
		m2, _ := eng.d.LoadManifest(file2)
		if m2 == nil {
			t.Errorf("expected file2 manifest to survive")
		}

		// Delete the second file
		eng.DeleteFile(file2)

		// Block must be wiped
		if len(eng.it.Indexes) != 0 {
			t.Errorf("expected IndexTable to be empty after deleting both files")
		}
	})
}

func TestReadBlock(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		contentSize int
		expectedEOF bool
	}{
		{
			name:        "Scenario 1: File smaller than the block (Padding required)",
			contentSize: 10,
			expectedEOF: true,
		},
		{
			name:        "Scenario 2: File exactly the block size",
			contentSize: block_size,
			expectedEOF: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			// ACT
			err = readBlock(f, buf, &endOfFile)

			// ASSERT
			if err != nil && !errors.Is(err, io.EOF) {
				t.Fatalf("unexpected error in readBlock: %v", err)
			}

			if endOfFile != tt.expectedEOF {
				t.Errorf("expected endOfFile = %v, found %v", tt.expectedEOF, endOfFile)
			}

			// Verify padding
			if tt.contentSize < block_size {
				if buf[tt.contentSize] != 0 {
					t.Errorf("padding failed! the byte at position %d is not zero", tt.contentSize)
				}
				if buf[block_size-1] != 0 {
					t.Errorf("padding failed! the last byte of the buffer is not zero")
				}
			}
		})
	}
}
