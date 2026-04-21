package engine

import (
	"errors"
	"io"
	"log"
	"os"
)

// Engine holds the engine state, including the volume manager,
// index table, and manifest index.
type Engine struct {
	vm VolumeManager
	it IndexTable
	mi ManifestIndex
}

// NewEngine creates a new Engine with the given volume manager, index table, and manifest index.
func NewEngine(vm VolumeManager, it IndexTable, mi ManifestIndex) *Engine {
	return &Engine{
		vm: vm,
		it: it,
		mi: mi,
	}
}

const (
	blockSize int64 = 4096
)

// processFile will open the file, read it in chunks of 4096 bytes,
// and process each block.
func (e *Engine) ProcessFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		log.Println("error opening file:", file)
		return err
	}

	// buf will hold the block data.
	// It will read 4096 bytes at a time.
	buf := make([]byte, blockSize)

	// offset will hold the current offset in the file.
	// It will be incremented after each iteration.
	offset := 0

	// Once we have the buffer, we start chunking the
	// file and processing its parts.
loop:
	for {
		// Reading the chunks from the file
		// and storing it in buf.
		_, err := f.Read(buf)
		// Checks if the error received happened because
		// the file had ended.
		if errors.Is(err, io.EOF) {
			log.Println("file has ended")
			break loop
		} else if err != nil {
			log.Println("error ocurred reading file:", file)
			return err
		}

		// Once the bytes are in buf, we
		// start processing the block.

		// At the end of the iteration, increase the offset
		// and move the reading area of the buf.
		offset++
		f.Seek(blockSize*int64(offset), 0)
	}

	return nil
}

// processBlock will get the block hash,
// check for duplicity in the index table
// and decide wheter to store it or create
// update the refCount.
func (e *Engine) processBlock(block []byte, offset int) error {
	// Generating the block hash.
	hash, err := generateHash(block)
	if err != nil {
		log.Println("error hashing block.")
		return err
	}

	// Checking if the hash already exists
	// in the index table.
	if ok := e.it.checkHash(hash); !ok {
		// If the hash doesn't exist, create a new index
		// and add it to the index table.
		i := NewIndex(hash, offset)
		e.it.addIndex(i)
	}

	return nil
}
