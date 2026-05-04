package volume

import (
	"fmt"
	"io"
	"os"
)

// File represents a volume file on the filesystem.
type File struct {
	file *os.File
	info os.FileInfo
}

// newFile creates a new File instance by opening the volume file on the filesystem.
func newFile(filePath string) (*File, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
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
	// Seek back to the start of the file before reading.
	f.file.Seek(0, io.SeekStart)

	buf := make([]byte, BitMapSize)
	// Read the bitmap from the file at offset 0.
	_, err := f.file.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading bitmap: %w", err)
	}

	return buf, nil
}

// writeBitMap writes the bitmap to the volume file.
func (f *File) writeBitMap(bitMap []byte) error {
	_, err := f.file.Write(bitMap)
	if err != nil {
		return fmt.Errorf("error writing bitmap: %w", err)
	}

	return nil
}

// Auxiliary function to sotre a byte at a given position in the file.
func (f *File) writeBlock(block []byte, position int) (error) {
	// Calculating the absolute offset.
	offset := int64(BitMapSize + (position * 8))

	// Writing the block to the disk.
	_, err := f.file.WriteAt(block, offset)
	if err != nil {
		return fmt.Errorf("error writing block to file: %w", err)
	}

	return nil
}