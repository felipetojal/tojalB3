package volume

import (
	"log"
)

// VolumeManager manages the volume file and bit map.
type VolumeManager struct {
	// volumeFile is the file that holds the raw bytes
	// of the stored files.
	volumeFile *File
	// bitMap represents the occupied and free positions
	// in the volumeFile.
	bitMap *BitMap
}

// NewVolumeManager creates a new VolumeManager instance.
func NewVolumeManager() *VolumeManager {
	// Create a new volume file.
	volumeFile, err := newFile()
	if err != nil {
		log.Panic(err)
	}

	// Load the bit map from the volume file.
	bitMap, err := loadBitMap(volumeFile)
	if err != nil {
		log.Panic(err)
	}

	// Initialize the volume manager with the volume file and bit map.
	return &VolumeManager{
		volumeFile: volumeFile,
		bitMap:     bitMap,
	}
}
