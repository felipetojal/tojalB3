package volume

const (
	BitMapSize = 4096
)

// BitMap represents the occupied and free positions
// in the volume file.
//
// It is stored in the volume.dat file as a header
// (bytes 0 - 4095).
type BitMap struct {
	// 1 - Occupied
	// 0 - Free
	bitMap []byte
}

// occupyPosition occupies the position at the given index.
func (b *BitMap) occupyPosition(position int) {
	b.bitMap[position] = 1
}

// freePosition frees the position at the given index.
func (b *BitMap) freePosition(position int) {
	b.bitMap[position] = 0
}

// loadBitMap loads the bit map from the volume file.
func loadBitMap(volume *File) (*BitMap, error) {
	bitmap := make([]byte, BitMapSize)
	// If the volume file is empty, initialize it with the initial file size.
	if volume.info.Size() == 0 {
		return &BitMap{
			bitMap: bitmap,
		}, nil
	}

	// Read the bitmap from the volume file.
	bitmap, err := volume.readBitMap()
	if err != nil {
		return nil, err
	}

	return &BitMap{
		bitMap: bitmap,
	}, nil
}
