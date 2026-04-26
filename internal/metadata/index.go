package metadata

import "fmt"

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
		hash:     hash,
		address:  address,
		refCount: 1,
	}
}

// IndexTable is a struct that maps block hashes
// to Index objects. This will allow O(1) lookups.
type IndexTable struct {
	indexes map[string]Index
}

// newIndexTable creates a new IndexTable.
func newIndexTable() *IndexTable {
	return &IndexTable{
		indexes: make(map[string]Index),
	}
}

// getIndex receives a block hash and returns the
// index with the block information.
func (it *IndexTable) getIndex(hash string) (*Index, error) {
	if !checkExistence(it, hash) {
		return nil, fmt.Errorf("error: hash does not exist in IndexTable.")
	}

	index := it.indexes[hash]

	return &index, nil
}

// addIndex receives an Index and adds it to the IndexTable.
// If the index does not yet exists, we must cerate a new
// key-value pair. If it already exists and we are trying to
// add it again, we must increase the refCount.
func (it *IndexTable) addIndex(i Index) {
	// Getting the original block hash.
	hash := i.hash
	// Declaring an Index variable in the function scope with
	// the same value as i.
	index := i
	// Check the existence of the given hash.
	if checkExistence(it, hash) {
		// If it the hash does exist, we update the refCount.
		index = it.indexes[hash]
		index.refCount = index.refCount + 1
	}

	// Storing the Index in the map
	it.indexes[hash] = index
}

// checkExistence is an auxiliary function used to check
// if a block hash already exists in the IndexTable.
func checkExistence(it *IndexTable, hash string) bool {
	_, ok := it.indexes[hash]
	return ok
}
