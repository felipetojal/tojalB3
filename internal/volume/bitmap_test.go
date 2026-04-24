package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFreePosition tests the freePosition method of the BitMap struct.
func TestFreePosition(t *testing.T) {
	b := &BitMap{
		bitMap: make([]byte, (BitMapSize / 10)),
	}

	// Slice of valid indexes to test occupyPosition.
	testValidIndexes := []int{1000, 19, 1014, 89, 294, 755}

	// Test occupyPosition with valid indexes.
	for _, value := range testValidIndexes {
		// Test occupyPosition with valid indexes.
		err := b.occupyPosition(value)
		assert.Nil(t, err)

		// Test freePosition with valid indexes.
		err = b.freePosition(value)
		assert.Nil(t, err)

		// Test assertBitIsZero with valid indexes.
		assertBitIsZero(t, b, value)
	}

	// Test occupyPosition with invalid indexes.
	invalidIndexes := []int{124943, -34, -2, 213484, 91340}
	for _, value := range invalidIndexes {
		err := b.freePosition(value)
		assert.NotNil(t, err)
	}
}

// TestOccupyPosition tests the occupyPosition method of the BitMap struct.
func TestOccupyPosition(t *testing.T) {
	b := &BitMap{
		bitMap: make([]byte, (BitMapSize / 10)),
	}

	// Slice of valid indexes to test occupyPosition.
	testValidIndexes := []int{10, 12, 502, 13, 94, 55}

	// Iterating over the valid indexes and making
	// sure the bit is set to 1 after occupyPosition is called.
	for _, value := range testValidIndexes {
		// First we need to occupy the bits.
		err := b.occupyPosition(value)
		assert.Nil(t, err)

		// Test assertBitIsSet with valid indexes.
		assertBitIsSet(t, b, value)
	}

	// Iterating over the invalid indexes and making
	// sure an error is returned from the access.
	testInvalidIndexes := []int{1000000, 23459000, -4, -2, 59250000, 1000000}
	for _, value := range testInvalidIndexes {
		err := b.occupyPosition(value)
		assert.Error(t, err)
	}
}

// assertBitIsSet checks if the bit at the given position is set to 1.
func assertBitIsSet(t *testing.T, b *BitMap, position int) {
	t.Helper()
	assert := assert.New(t)

	bitMapIndex, byteIndex := calculateIndex(t, position)

	// Checking if the bit at the given bitMapIndex and byteIndex is set to 1.
	// This operation will isolate the bit we want to look at.
	//
	// Say b.bitMap[bitMapIndex] = 00010101 and byteIndex = 3, we will have
	// 00010101 & (00000100) = 00000100 = 4.
	//
	// If the value is different from 0, it is valid.
	bit := (b.bitMap[bitMapIndex] & (1 << byteIndex))

	assert.NotEqual(0, bit)
}

// assertBitIsZero checks if the bit at the given position is set to 0.
func assertBitIsZero(t *testing.T, b *BitMap, position int) {
	t.Helper()
	assert := assert.New(t)

	bitMapIndex, byteIndex := calculateIndex(t, position)

	// Checking if the bit at the given bitMapIndex and byteIndex is set to 1.
	// This operation will isolate the bit we want to look at.
	//
	// Say b.bitMap[bitMapIndex] = 00010101 and byteIndex = 3, we will have
	// 00010101 & (00000100) = 00000100 = 4.
	//
	// If the value is different from 0, it is valid.
	bit := (b.bitMap[bitMapIndex] & (1 << byteIndex))

	assert.Equal(byte(0), bit)
}

// calculateIndex calculates the bitMapIndex and byteIndex for a given position.
func calculateIndex(t *testing.T, position int) (int, int) {
	t.Helper()

	bitMapIndex := position / 8
	byteIndex := position % 8
	return bitMapIndex, byteIndex
}
