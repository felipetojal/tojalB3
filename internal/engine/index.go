package engine

// Index represents a block index in the engine.
type Index struct {
	// The SHA-256 hash associated with the block bytes.
	hash string
	// The address within the volume.dat file.
	address int
	// The number of times that block is stored (avoiding
	// deduplication).
	refCount int
}

// NewIndex creates a new Index with the given hash and address.
func NewIndex(hash string, address int) Index {
	return Index{
		hash:     hash,
		address:  address,
		refCount: 1,
	}
}

// IxTable defines the interface for the index table.
type IxTable interface {
	checkHash(string) bool
	inserIndex(Index) error
}

// IndexTable holds a map of block hashes to their corresponding Index.
type IndexTable struct {
	// Allows O(1) lookups.
	indexes map[string]Index
}

// checkHash will check for block duplicity.
func (i *IndexTable) checkHash(hash string) bool {
	_, ok := i.indexes[hash]
	return ok
}

func (i *IndexTable) addIndex(index Index) {
	i.indexes[index.hash] = index
}
