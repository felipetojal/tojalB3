package volume

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewFile tests the creation of a new file.
func TestNewFile(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory for the test.
	testFile := testFilePath(t, "volume_test.dat")

	// Create a new file.
	f, err := newFile(testFile)
	assert.Nil(err)
	assertFileNotNil(t, f)
}

// TestReadBitMap tests reading the bitmap from a file.
func TestReadBitMap(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory for the test.
	testFile := testFilePath(t, "volume_test.dat")

	// Create a new file.
	f, err := newFile(testFile)
	assert.Nil(err)
	assertFileNotNil(t, f)

	// Write the bit map to the file.
	bitMap := []byte{34, 132, 54, 32, 9}
	err = f.writeBitMap(bitMap)
	assert.Nil(err)

	// Read the bit map from the file.
	b, err := f.readBitMap()
	assert.Nil(err)
	assert.NotEqual(0, len(b))
}

// TestWriteBitMap tests writing the bitmap to a file.
func TestWriteBitMap(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory for the test.
	testFile := testFilePath(t, "volume_test.dat")

	// Create a new file.
	f, err := newFile(testFile)
	assert.Nil(err)
	assertFileNotNil(t, f)

	// Write the bit map to the file.
	bitMap := []byte{34, 132, 54, 32, 9}
	err = f.writeBitMap(bitMap)
	assert.Nil(err)
}

func TestWriteBlock(t *testing.T) {
	testFile := testFilePath(t, "volume_test.dat")

	f, err := newFile(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	block := bytes.Repeat([]byte("1"), block_size)
	assert.NotNil(t, block)
	err = f.writeBlock(block, 89)
	assert.NoError(t, err)
}

func TestDeleteBlock(t *testing.T) {
	testFile := testFilePath(t, "volume_test.dat")

	f, err := newFile(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	block := bytes.Repeat([]byte("1"), block_size)
	assert.NotNil(t, block)
	err = f.writeBlock(block, 89)
	assert.NoError(t, err)

	err = f.deleteBlock(89)
	assert.NoError(t, err)
}

func TestReadBlock(t *testing.T) {
	testFile := testFilePath(t, "volume_test.dat")

	f, err := newFile(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)

	block := bytes.Repeat([]byte("1"), block_size)
	assert.NotNil(t, block)
	err = f.writeBlock(block, 89)
	assert.NoError(t, err)

	blockRead, err := f.readBlock(89)
	assert.NoError(t, err)
	assert.Equal(t, block, blockRead)
}

// testFilePath returns the path to a test file in the temporary directory.
func testFilePath(t *testing.T, filePath string) string {
	t.Helper()
	tempDir := t.TempDir()
	testFilePath := tempDir + filePath
	return testFilePath
}

// assertFileNotNil asserts that the file was created successfully.
func assertFileNotNil(t *testing.T, f *File) {
	assert := assert.New(t)
	t.Helper()

	// Assert that the file was created successfully.
	assert.NotNil(f)
	assert.NotNil(f.file)
	assert.NotNil(f.info)
}
