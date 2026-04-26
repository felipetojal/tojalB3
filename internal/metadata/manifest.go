package metadata

import "time"

// Manifest is a struct that holds
// file metadata.
type Manifest struct {
	// File basic information.
	fileName string
	fileDir  string
	fileSize int64

	// All the block hashes associated to the file.
	blocks []string

	// Timestamp of creation
	storedAt time.Time
}

// newManifest creates a new Manifest.
func newManifest(fileName, fileDir string, fileSize int64, blocks []string) *Manifest {
	return &Manifest{
		fileName: fileName,
		fileDir:  fileDir,
		fileSize: fileSize,
		blocks:   blocks,
		storedAt: time.Now().UTC(),
	}
}

// ManifestTable is a struct that holds all 
// the Manifests saved in the database.
type ManifestTable struct {
	manifests map[string]Manifest
}

// newManifestTable creates a new ManifestTable.
func newManifestTable() (*ManifestTable) {
	return &ManifestTable{
		manifests: make(map[string]Manifest),
	}
}