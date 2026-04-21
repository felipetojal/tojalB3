package storage

import (
	"os"
	"time"
)

// Manifest holds the metadata of a file stored in the database.
type Manifest struct {
	fileName string
	filePath string
	fileSize uint64

	// Each manifest will have a slice
	// of the block hashes that belong
	// to the original file
	blocks []string

	// Storing the creation timestamp
	// for logging purposes.
	storedAt time.Time
}

// NewManifest creates a new Manifest instance with the given parameters.
func NewManifest(fileName, filepath string, fileSize uint64, blocks []string) Manifest {
	return Manifest{
		fileName: fileName,
		filePath: filepath,
		fileSize: fileSize,
		blocks:   blocks,
		storedAt: time.Now().UTC(),
	}
}

// Each index represent a block. Since the hash field
// is unique, if another block has the same hash field,
// we don´t store it separately, we increase the refCount
// field to represent duplicity.
type Index struct {
	// The block hash associated with this index.
	hash string

	// Represents the address within the volume.dat file.
	address uint64

	// Counting the number of times this block is referenced.
	refCount int
}

// NewIndex creates a new Index instance with the given parameters.
func NewIndex(hash string, address uint64, refCount int) Index {
	return Index{
		hash:     hash,
		address:  address,
		refCount: refCount,
	}
}

// VolumeManager handles the volume.dat file and provides methods for block storage.
type VolumeManager struct {
	file   *os.File
	bitMap []uint8
}
