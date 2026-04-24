package volume

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFreePosition tests the freePosition method of the BitMap struct.
func TestFreePosition(t *testing.T) {
	assert.New(t)

	b := &BitMap{
		bitMap: make([]byte, (BitMapSize / 10)),
	}

	testIndexes := []int{1000, 19, 10144, 89, 294, 755}

	for index := range testIndexes {
		b.occupyPosition(index)
	}

}

// TestOccupyPosition tests the occupyPosition method of the BitMap struct.
func TestOccupyPosition(t *testing.T) {
	assert.New(t)

	b := &BitMap{
		bitMap: make([]byte, (BitMapSize / 10)),
	}

	// Iterating over the valid indexes and making
	// sure no error is returned from the access.
	testValidIndexes := []int{10, 12, 502, 13, 94, 55}
	for _, value := range testValidIndexes {
		err := b.occupyPosition(value)
		assert.Nil(t, err)
	}

	// calcIndexes is a helper function that calculates the bitMapIndex and byteIndex
	// for a given position in the bitMap.
	calcIndexes := func(position int) (int, int) {
		bitMapIndex := position / 8
		byteIndex := position % 8
		return bitMapIndex, byteIndex
	}

	// equalTo1 is a helper function that checks if the bit at the given
	// bitMapIndex and byteIndex is set to 1.
	equalTo1 := func(bitMapIndex, byteIndex int) error {
		bit := b.bitMap[bitMapIndex] << byteIndex
		if bit != 1 {
			return fmt.Errorf("bit at bitMapIndex %d and byteIndex %d is not 1, is equal to %v",
				bitMapIndex, byteIndex, bit)
		}
		return nil
	}

	// Iterating over the valid indexes and making
	// sure the bit is set to 1 after occupyPosition is called.
	for _, value := range testValidIndexes {
		bitMapIndex, byteIndex := calcIndexes(value)
		err := equalTo1(bitMapIndex, byteIndex)
		assert.Nil(t, err)
	}

	// Iterating over the invalid indexes and making
	// sure an error is returned from the access.
	testInvalidIndexes := []int{1000000, 23459000, -4, -2, 59250000, 1000000}
	for index := range testInvalidIndexes {
		err := b.occupyPosition(index)
		assert.Error(t, err)
	}
}
