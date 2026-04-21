package engine

import (
	"os"
	"time"
)

// ManifestIndex holds a map of file manifests indexed by file name hash.
// It will be stored in a key-value database.
type ManifestIndex struct {
	manifest map[string]Manifest
}

// Manifest holds the metadata for a file.
type Manifest struct {
	fileName string
	filePath string
	fileSize uint64

	// All block hashes belonging to the file
	// are stored here.
	blocks map[string]Index

	storedAt time.Time
}

// NewManifest creates a new Manifest instance.
func NewManifest(fileName, filePath string, fileSize uint64, blocks map[string]Index) Manifest {
	return Manifest{
		fileName: fileName,
		filePath: filePath,
		fileSize: fileSize,
		blocks:   blocks,
		storedAt: time.Now().UTC(),
	}
}

// VolumeManager manages the volume file and bit map.
type VolumeManager struct {
	// volumeFile is the file that holds the raw bytes
	// of the stored files.
	volumeFile *os.File
	// bitMap represents the occupied and free positions
	// in the volumeFile.
	bitMap []int8
}
