package volume

import (
	"fmt"
	"os"
)

// File represents a volume file on the filesystem.
type File struct {
	file *os.File
	info os.FileInfo
}

type file interface {
	readBitMap() ([]byte, error)
}

// newFile creates a new File instance by opening the volume file on the filesystem.
func newFile() (*File, error) {
	f, err := os.OpenFile("/tmp/data/volume.dat", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("error generating stats: %w", err)
	}

	return &File{
		file: f,
		info: info,
	}, nil
}

// readBitMap reads the bitmap from the volume file.
func (f *File) readBitMap() ([]byte, error) {
	buf := make([]byte, BitMapSize)
	_, err := f.file.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading bitmap: %w", err)
	}

	return buf, nil
}
