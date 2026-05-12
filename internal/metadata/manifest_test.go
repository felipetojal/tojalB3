package metadata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewManifest is responsible for testing
// the Manifest constructor.
func TestNewManifest(t *testing.T) {
	a := assert.New(t)

	fileName := "toddynho.png"
	fileSize := int64(100)
	blocks := []string{"oi", "ida", "volta", "naruto"}
	m := NewManifest(fileName, fileSize)
	m.AddBlock(blocks...)
	a.NotNil(m)
}

func TestAddBlock(t *testing.T) {
	tempDir := t.TempDir()

	fileName := "jojo.txt"
	filePath := filepath.Join(tempDir, fileName)
	f, err := os.Create(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	m := NewManifest(fileName, 7000)
	assert.NotNil(t, m)

	blocks := []string{"ogwrio", "novn249gn249"}

	m.AddBlock(blocks...)

	l := len(m.Blocks)
	assert.Equal(t, 2, l)
}
