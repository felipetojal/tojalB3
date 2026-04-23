package main

import (
	"fmt"
	"log"
	"os"

	"github.com/felipetojal/tojalB3/internal/engine"
)

const block_size = 4096

func main() {
	f, err := os.Open("tojalB3-Diagrams/tojalB3-ArchitecturalDiagramImage.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	buf := make([]byte, block_size)
	offset := 0
loop:
	for {
		n, err := f.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		readBuf := buf[:n]

		hash, err := engine.GenerateHash(readBuf)
		if err != nil {
			log.Println("error ocurred generating hash.")
			return
		}

		log.Printf("BLOCK: %d\n", offset)
		log.Printf("LENGTH: %v\n", len(readBuf))
		log.Printf("HASH: %v\n", []byte(hash))
		log.Printf("CONTENT: %v\n", []byte(readBuf))

		// If the number of bytes read is less than
		// the block size, we assume that it is the last block.
		if n < int(block_size) {
			break loop
		}

		offset++
		f.Seek(int64(block_size*offset), 0)
	}
}

func run() error {
	// If the file already exists, it will be opened,
	// if not, it will be created.
	_, err := os.OpenFile("/data/volume.dat", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("error opening volume.dat file: %w", err)
	}

}
