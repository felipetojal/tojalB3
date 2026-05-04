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
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %v: %w", filePath, err)
	}
	fileInfo, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error creating fileInfo: %w", err)
	}
	fileSize := fileInfo.Size()
	defer f.Close()

	// Since the file does not exist, we create a new Manifest.
	mani := metadata.NewManifest(filePath, fileSize)

	// Creating and reading the buffer.
	buf := make([]byte, block_size)
	endOfFile := false
	offset := 0
	// Now, we need to read each chunk of the file.
	for !endOfFile {
		// Reading the block and checking possible errors.
		err := readBlock(f, buf)
		if !errors.Is(err, io.EOF) {
			return err
		}
		// If the end of the file was reached, then
		// we must signal so stop the loop on this
		// iteration.
		if errors.Is(err, io.EOF) {
			endOfFile = true
		}

		// Now that the bytes are in the buffer, we must
		// create the blocks (indexes).
		blockHash, err := hash.GenerateHash(buf)
		if err != nil {
			return err
		}
		// Adding the block hash to the manifest pointer.
		mani.AddBlock(blockHash)

		// Declaring the index.
		var i *metadata.Index

		// Before creating a new index, we check if that same
		// block already exists in the volume.
		v, ok := e.it.Indexes[blockHash]
		if ok {
			v.RefCount++
			e.it.Indexes[blockHash] = v
		} else {
			i = metadata.NewIndex(blockHash, offset)
		}
		
		// At the end of the iteration, we must move
		// the file pointer to the next chunk of bytes.
		offset++
		of, err := f.Seek(0, offset*block_size)
		if err != nil {
			return fmt.Errorf("error changing offset: %w", err)
		}
	}

	return nil
}

// readBlock is responsible for reading the block from
// the file. It returns the buffer from the read and
// an error.
func readBlock(f *os.File, buf []byte) error {
	n, err := f.Read(buf)
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading block: %w", err)
	}

	// If less bytes than what can be stored in the block were
	// read, that means we arrived at the last block of the file.
	// So we need to complete the remaining bytes with zeros(padding).
	if n < block_size {
		clear(buf[n:])
		return err
	}

	return nil
}
