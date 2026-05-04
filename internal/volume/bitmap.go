package volume

import (
	"fmt"
)

const (
	BitMapSize = 40960
	block_size = 4096
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

// freePosition frees the position at the given index.
// position is the block index in the volume.
func (b *BitMap) freePosition(position int) error {
	if (position > BitMapSize) || (position < 0) {
		return fmt.Errorf("index %d out of bounds", position)
	}
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

	return nil
}

// occupyPosition occupies the position at the given index.
// position is the block index in the volume.
func (b *BitMap) occupyPosition(position int) error {
	if (position > BitMapSize) || (position < 0) {
		return fmt.Errorf("index %d out of bounds", position)
	}

	// Getting the right indexes.
	bitMapIndex := position / 8
	byteIndex := position % 8

	// Setting the byteIndex to 1, indicating that it is
	// now occupied.
	b.bitMap[bitMapIndex] = b.bitMap[bitMapIndex] | (1 << byteIndex)

	return nil
}

// getNextFreePosition returns the position of the next free block.
func (b *BitMap) getNextFreePosition() (int, error) {
	// Iterating over the bitmap.
	for bitMapIndex := 0; bitMapIndex < BitMapSize; bitMapIndex++ {

		// The maximum value a byte can have is 255. In this
		// case, it would be full.
		if b.bitMap[bitMapIndex] != byte(255) {

			// Searching which bit is free.
			bitIndex := getFreeBit(b.bitMap[bitMapIndex])

			absolutePosition := (bitMapIndex * 8) + bitIndex

			return absolutePosition, nil
		}
	}

	return -1, fmt.Errorf("volume is full: no free blocks available")
}

// getFreeBit analyzes a byte and returns the bit free position (0 a 7).
func getFreeBit(b byte) int {
	for index := 0; index < 8; index++ {
		if (b & (1 << index)) == 0 {
			return index // Returns only the bit number (ex: 3).
		}
	}
	return -1
}
