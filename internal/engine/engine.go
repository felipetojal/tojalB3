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

	// Opening the file.
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
	// Now, we need to read each chunk of the file.
	for !endOfFile {
		// Reading the block and checking possible errors.
		err := readBlock(f, buf, &endOfFile)
		if !errors.Is(err, io.EOF) && err != nil {
			return err
		}

		// Now that the bytes are in the buffer, we must
		// create the blocks (indexes).
		blockHash, err := hash.GenerateHash(buf)
		if err != nil {
			return err
		}

		mani.AddBlock(blockHash)

		// Validating the block existence.
		if exist := e.it.CheckExistence(blockHash); exist {
			index := e.it.Indexes[blockHash]
			index.RefCount++
			e.it.Indexes[blockHash] = index
		} else {
			// If the block is not found in the index table,
			// we must add it and create the index.
			if err := e.storeBlock(buf, blockHash); err != nil {
				return err
			}
		}
	}

	// After the operation is done, we save the manifest in
	// the database.
	if err := e.d.StoreManifest(mani); err != nil {
		return fmt.Errorf("error storing manifest: %w", err)
	}

	// Storing the index table in the database.
	if err := e.d.StoreIndexTable(e.it); err != nil {
		return fmt.Errorf("error storing index table: %w", err)
	}

	return nil
}

// Auxiliary funtion to encapsulate the logic to store a block
// in the volume file.
func (e *Engine) storeBlock(block []byte, hash string) error {
	// Checking if the hash already exists
	// Adding the block hash to the manifest pointer.
	pos, err := e.v.StoreBlock(block)
	if err != nil {
		return err
	}
	e.it.Indexes[hash] = *metadata.NewIndex(hash, pos)

	return nil
}

// readBlock is responsible for reading the block from
// the file. It returns the buffer from the read and
// an error.
func readBlock(f *os.File, buf []byte, endOfFile *bool) error {
	n, err := f.Read(buf)
	if !errors.Is(err, io.EOF) && err != nil {
		return fmt.Errorf("error reading block: %w", err)
	}

	// If less bytes than what can be stored in the block were
	// read, that means we arrived at the last block of the file.
	// So we need to complete the remaining bytes with zeros(padding).
	if n < block_size {
		clear(buf[n:])
		*endOfFile = true
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
