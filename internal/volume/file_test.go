package volume

import (
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

	b, err := f.readBitMap()
	assert.Nil(err)
	assert.NotEqual(0, len(b))
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
