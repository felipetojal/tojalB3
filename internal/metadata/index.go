package metadata

// Index is a struct that holds data
// of a certain block.
type Index struct {
	// hash is made based on the block bytes.
	hash string

	// address is the position of the block
	// in the volume file. This helps for lookups.
	address int

	// refCount is responsible for avoiding
	// deduplication. If a block is being processed
	// and its hash is equal to another block already
	// stored in memory, then they are equal, hence
	// there is no need to store it again.
	refCount int
}

// newIndex creates a new Index.
func newIndex(hash string, address int) *Index {
	return &Index{
		hash:    hash,
		address: address,
	}
}
