package engine

import (
	"fmt"
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
		n, err := f.Read(buf)
		if err != nil {
			log.Println("error reading bytes from file")
			return fmt.Errorf("error reading bytes: %w", err)
		}

	}

	return nil
}
