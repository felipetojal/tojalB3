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
func newManifest(fileName string, fileSize int64, blocks []string) *Manifest {
	return &Manifest{
		FileName: fileName,
		FileSize: fileSize,
		Blocks:   blocks,
		StoredAt: time.Now().UTC(),
	}
}

func (m *Manifest) prefix() string {
	return "mani:"
}
