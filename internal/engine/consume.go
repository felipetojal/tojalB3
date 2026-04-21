package engine

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	BLOCK_SIZE int64 = 4096
)

// processFile holds the main logic of the operation.
// It will open a file, iterate over it till its over
// and call the proper functions to handle the byte
// storage.
func processFile(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Creating the buffer with the same size
	// as the block size. This way, every read
	// from the program will be a separate block.
	buf := make([]byte, BLOCK_SIZE)
	for {
		// Reading the fixed block size from the file.
		n, err := f.Read(buf)
		if err != nil {
			log.Println("error reading bytes from file")
			return fmt.Errorf("error reading bytes: %w", err)
		}

		// This will verify if the error reading the file
		// occurred because of the end of the file.
		if errors.Is(err, io.EOF) {
			log.Println("file has reached the end.")
			break
		}

		// Generating the hash based on SHA-256.
		hash, err := generateHash(buf)
		if err != nil {
			log.Println("error generating hash")
			return err
		}

		/*
		* Once we have the block hash, we must verify if it
		* already exists in the index table. If it already exists,
		* we increment the reference count. If it doesn´t exist,
		* we must find a place to store it using the bit map.
		 */

		// Skipping the block size to read the next block.
		f.Seek(BLOCK_SIZE, 1)

	}

	return nil
}
