package engine

import (
	"errors"
	"fmt"
	"io"
	"os"

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
	// Before openning the file, we must check for its existence.

	// Openning the file.
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %v: %w", filePath, err)
	}
	defer f.Close()

	// Now, we need to read each chunk of the file.
	for {
		// Creating and reading the buffer.
		buf := make([]byte, block_size)
		// Reading the block using an auxiliary function.
		err := readBlock(f, buf)
		if !errors.Is(err, io.EOF) {
			return err
		}
		
	}

	return nil
}

// readBlock is responsible for reading the block from
// the file. It returns the buffer from the read and
// an error.
func readBlock(f *os.File, buf []byte) (error) {
	n, err := f.Read(buf)
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading block: %w", err)
	}

	// If less bytes than what can be stored in the block were
	// read, that means we arrived at the last block of the file.
	// So we need to complete the remaining bytes with zeros(padding)
	if n < block_size {
		clear(buf[n:])
		return err
	}

	return nil
}
