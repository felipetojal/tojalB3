package engine

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/felipetojal/tojalB3/internal/hash"
	"github.com/felipetojal/tojalB3/internal/metadata"
	"github.com/felipetojal/tojalB3/internal/volume"
)

// Engine is responsible for handling the
// main logic of the program.
type Engine struct {
	v  *volume.VolumeManager
	d  *metadata.Database
	it *metadata.IndexTable
}

/*
What is Engine responsible for?

1 - It must support four operations: store file, load file, delete file, and list all files.
2 - When storing the file, it must convert it into chunks and store it in the volume.
*/

const block_size = 4096

// StoreFile will receive the filePath, open the file, split it into chunks,
// and store each chunk in the volume.
func (e *Engine) StoreFile(filePath string) error {
	// Now we must check if the file already exists in the database.
	m, err := e.d.LoadManifest(filePath)
	if m != nil {
		return fmt.Errorf("error file already exists: %w", err)
	}

	// Openning the file.
	f, fileSize, err := openFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Since the file does not exist, we create a new Manifest.
	mani := metadata.NewManifest(filePath, fileSize)

	// Creating and reading the buffer.
	buf := make([]byte, block_size)
	endOfFile := false
	fileOffset := 0
	// Now, we need to read each chunk of the file.
	for !endOfFile {
		// Reading the block and checking possible errors.
		err := readBlock(f, buf, endOfFile)
		if !errors.Is(err, io.EOF) {
			return err
		}

		// Now that the bytes are in the buffer, we must
		// create the blocks (indexes).
		blockHash, err := hash.GenerateHash(buf)
		if err != nil {
			return err
		}

		// Checking if the hash already exists
		// Adding the block hash to the manifest pointer.
		mani.AddBlock(blockHash)

		// At the end of the iteration, we must move
		// the file pointer to the next chunk of bytes.
		if err := advanceOffset(f, fileOffset); err != nil {
			return err
		}
	}

	return nil
}

// readBlock is responsible for reading the block from
// the file. It returns the buffer from the read and
// an error.
func readBlock(f *os.File, buf []byte, endOfFile bool) error {
	n, err := f.Read(buf)
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading block: %w", err)
	}

	// If less bytes than what can be stored in the block were
	// read, that means we arrived at the last block of the file.
	// So we need to complete the remaining bytes with zeros(padding).
	if n < block_size {
		clear(buf[n:])
		endOfFile = true
		return err
	}

	return nil
}

// Auxiliary function to open a given filePath.
func openFile(filePath string) (*os.File, int64, error) {
	// Openning the file.
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, 0, fmt.Errorf("error opening file %v: %w", filePath, err)
	}
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("error creating fileInfo: %w", err)
	}
	fileSize := fileInfo.Size()

	return f, fileSize, nil
}

// Auxiliary function to advance the pointer of the
// file offset.
func advanceOffset(f *os.File, fileOffset int) error {
	_, err := f.Seek(int64(fileOffset)*block_size, 0)
	fileOffset++
	if err != nil {
		return fmt.Errorf("error changing offset: %w", err)
	}
	return nil
}
