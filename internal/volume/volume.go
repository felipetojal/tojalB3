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
func NewVolumeManager(filePath string) *VolumeManager {
	// Create a new volume file.
	volumeFile, err := newFile(filePath)
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

// loadBitMap loads the bit map from the volume file.
func loadBitMap(f *File) (*BitMap, error) {
	bitmap := make([]byte, BitMapSize)
	// If the volume file is empty, initialize it with the initial file size.
	if f.info.Size() == 0 {
		// Initializing the bitMap with 0s.
		err := f.writeBitMap(bitmap)
		if err != nil {
			return nil, err
		}

		return &BitMap{
			bitMap: bitmap,
		}, nil
	}

	// Read the bitmap from the volume file.
	bitmap, err := f.readBitMap()
	if err != nil {
		return nil, err
	}

	return &BitMap{
		bitMap: bitmap,
	}, nil
}
