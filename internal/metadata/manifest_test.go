package metadata

import (
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
