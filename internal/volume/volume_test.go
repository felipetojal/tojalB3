package volume

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVolumeManager(t *testing.T) {
	assert := assert.New(t)

	tempDir := t.TempDir()
	vm := NewVolumeManager(tempDir + "volume_test.dat")
	assert.NotNil(vm)
}

// TestLoadBitMap tests the loadBitMap function
func TestLoadBitMap(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory and file for testing
	tempDir := t.TempDir()
	filePath := tempDir + "volume_test.dat"
	f, err := newFile(filePath)
	assert.Nil(err)
	assert.NotNil(f)
	defer f.file.Close()

	// Load the bit map from the file
	bitMap, err := loadBitMap(f)
	assert.Nil(err)
	assert.NotNil(bitMap)
}

func TestStoreBlock(t *testing.T) {
	block := bytes.Repeat([]byte("25954gn"), block_size)

	assert := assert.New(t)

	tempDir := t.TempDir()
	vm := NewVolumeManager(tempDir + "volume_test.dat")
	assert.NotNil(vm)

	_, err := vm.StoreBlock(block)
	assert.Nil(err)
}

func TestGetBlock(t *testing.T) {
	block := bytes.Repeat([]byte("1"), block_size)

	assert := assert.New(t)

	tempDir := t.TempDir()
	vm := NewVolumeManager(tempDir + "volume_test.dat")
	assert.NotNil(vm)

	pos, err := vm.StoreBlock(block)
	assert.Nil(err)

	b, err := vm.GetBlock(pos)
	assert.Nil(err)
	assert.Equal(b, block)
}

func TestVolumeDeleteBlock(t *testing.T) {
	block := bytes.Repeat([]byte("1"), block_size)

	assert := assert.New(t)

	tempDir := t.TempDir()
	vm := NewVolumeManager(tempDir + "volume_test.dat")
	assert.NotNil(vm)

	pos, err := vm.StoreBlock(block)
	assert.Nil(err)

	err = vm.DeleteBlock(pos)
	assert.Nil(err)
}
