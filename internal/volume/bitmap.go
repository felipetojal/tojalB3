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
	// Each byte contains 8 bits. Since we are only
	// storing 0s and 1s, we will split the bytes
	// to avoid memory waste.
	bitMap []byte
}

type bitMap interface {
	freePosition(int)
	occupyPosition(int)
	loadBitMap(*File) (*BitMap, error)
}

// freePosition frees the position at the given index.
// position is the block index in the volume.
func (b *BitMap) freePosition(position int) {
	/*
	* Example: Say the position (block index) is 45.
	* Since we are operating the bits of the bytes, and
	* each byte contains 8 bits, we make a mod operation
	* in the position to extract the right bitmap index.
	*
	* 43 / 8 = 5 -> bitmap index is 5
	* 43 % 8 = 3 -> location within the bitmap byte.
	 */

	// Getting the right indexes.
	bitMapIndex := position / 8
	byteIndex := position % 8

	// Setting the byteIndex to 0, indicating that it is
	// now free.
	b.bitMap[bitMapIndex] = b.bitMap[bitMapIndex] &^ (1 << byteIndex)
}

// occupyPosition occupies the position at the given index.
// position is the block index in the volume.
func (b *BitMap) occupyPosition(position int) {
	// Getting the right indexes.
	bitMapIndex := position / 8
	byteIndex := position % 8

	// Setting the byteIndex to 1, indicating that it is
	// now occupied.
	b.bitMap[bitMapIndex] = b.bitMap[bitMapIndex] | (1 << byteIndex)
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
