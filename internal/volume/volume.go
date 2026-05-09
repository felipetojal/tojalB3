package volume

import (
	"fmt"
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

func (v *VolumeManager) GetBlock(position int) ([]byte, error) {
	buf, err := v.volumeFile.readBlock(position)
	if err != nil {
		return nil, fmt.Errorf("error getting block: %w", err)
	}

	return buf, nil
}

// DeleteBlock receives the absolute position of the block
// and deletes it from the volume and frees the bitmap position.
func (v *VolumeManager) DeleteBlock(position int) error {
	if err := v.volumeFile.deleteBlock(position); err != nil {
		return err
	}

	if err := v.bitMap.freePosition(position); err != nil {
		return err
	}

	if err := v.volumeFile.writeBitMap(v.bitMap.bitMap); err != nil {
		return fmt.Errorf("error writing bit map: %w", err)
	}

	return nil
}

// This function is responsible for storing the block in
// the volume file and signaling the bitmap occupation.
// It returns the position of the block and a possible error.
func (v *VolumeManager) StoreBlock(block []byte) (int, error) {
	// Getting the position to store the block
	pos, err := v.bitMap.getNextFreePosition()
	if err != nil {
		return -1, fmt.Errorf("error getting freePosition storeBlock(): %w", err)
	}

	// Storing the block.
	if err := v.volumeFile.writeBlock(block, pos); err != nil {
		return -1, err
	}

	// Since the position was occupied, we must signal it
	// to the system.
	if err := v.bitMap.occupyPosition(pos); err != nil {
		return -1, fmt.Errorf("error occupyPosition storeBlock(): %w", err)
	}

	// Saving the new bitmap.
	if err := v.volumeFile.writeBitMap(v.bitMap.bitMap); err != nil {
		return -1, fmt.Errorf("error saving bitmap storeBlock: %w", err)
	}

	// pos is the absolue position of the block in the volume.
	return pos, nil
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
