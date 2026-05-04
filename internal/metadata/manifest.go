package metadata

import "time"

// Manifest is a struct that holds
// file metadata.
type Manifest struct {
	// File basic information.
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`

	// All the block hashes associated to the file.
	Blocks []string `json:"blocks"`

	// Timestamp of creation
	StoredAt time.Time `json:"storedAt"`
}

// newManifest creates a new Manifest.
func NewManifest(fileName string, fileSize int64) *Manifest {
	return &Manifest{
		FileName: fileName,
		FileSize: fileSize,
		Blocks:   make([]string, 0),
		StoredAt: time.Now().UTC(),
	}
}

// Auxiliary function to add blocks to the manifest.
func (m *Manifest) AddBlock(block ...string) {
	m.Blocks = append(m.Blocks, block...)
}

func (m *Manifest) prefix() string {
	return "mani:"
}
