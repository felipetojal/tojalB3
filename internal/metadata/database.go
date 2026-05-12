package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v4"
)

type Storable interface {
	prefix() string
}

// Database encapsulates the badger.DB type.
type Database struct {
	db    *badger.DB
	mutex sync.RWMutex
}

const (
	indexTableKey = "indexTable:"
	manifestKey   = "mani:"
)

// NewDatabase is the Database constructor. It takes in
// the directory path used to store the data and returns
// a pointer to Database and may return an error.
func NewDatabase(dirPath string) (*Database, error) {
	db, err := badger.Open(badger.DefaultOptions(dirPath))
	if err != nil {
		return nil, fmt.Errorf("unable to open badgerDB: %w", err)
	}

	return &Database{
		db:    db,
		mutex: sync.RWMutex{},
	}, nil
}

// ListFiles returns a list of all the files stored in the database.
func (d *Database) ListFiles() []string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	files := make([]string, 0)

	// Here we are creating a transaction
	// and an iterator that will go through
	// the badger items.
	txn := d.db.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	it := txn.NewIterator(opts)
	defer it.Close()

	// Iterating over the badger items and selecting
	// the manifests.
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.Key()
		kString := string(k)
		after, found := strings.CutPrefix(kString, manifestKey)
		if found {
			files = append(files, after)
		}
	}

	return files
}

// storeManifest is responsible for storing a Manifest in the
// database. It returns an error in case of failure
func (d *Database) StoreManifest(m *Manifest) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Converting the Manifest to bytes.
	mBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("error converting manifest to bytes: %w", err)
	}

	// Getting the Manifest key and converting it to bytes.
	key := []byte(manifestKey + m.FileName)

	// Creating the transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()
	err = txn.Set(key, mBytes)
	if err != nil {
		return fmt.Errorf("error storing key %v and value %v: %w", key, m, err)
	}

	// Committing and checking for errors.
	if err := txn.Commit(); err != nil {
		return fmt.Errorf("error commiting key %v and manifest %v: %w", key, m, err)
	}

	return nil
}

// loadManifest is responsible for loading a Manifest from
// the database. It receives the Manifest key as a parameter.
func (d *Database) LoadManifest(key string) (*Manifest, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	// Converting key from string to bytes.
	keyBytes := []byte(manifestKey + key)

	// Creating a transaction and loading object from
	// the database.
	txn := d.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("error loading key %v and Manifest %v from BadgerDB: %w", key, item, err)
	}

	// Once the object was retireved, we need to parse it
	// to the expected format.
	var m Manifest
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &m)
	})
	if err != nil {
		return nil, fmt.Errorf("error converting item to Manifest: %w", err)
	}

	// Returning the Manifest to the user.
	return &m, nil
}

// DeleteManifest will receive the file key (filePath) and
// delete the file in the database.
//
// If the file is not found or an error occur, an error will be returned.
// If the file is deleted successfully, error will be nil.
func (d *Database) DeleteManifest(key string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Creating the database transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()

	// Creating the ley and deleting the manifest.
	manifestKey := []byte(manifestKey + key)
	if err := txn.Delete(manifestKey); err != nil {
		return fmt.Errorf("error deleting manifest from database: %w", err)
	}

	// Commiting the operation.
	if err := txn.Commit(); err != nil {
		return fmt.Errorf("error commiting deletion from database: %w", err)
	}

	return nil
}

// StoreIndexTable is responsible for storing the IndexTable
// in the database. It returns an error.
func (d *Database) StoreIndexTable(it *IndexTable) error {
	// Converting the IndexTable to JSON bytes.
	itJson, err := json.Marshal(it)
	if err != nil {
		return fmt.Errorf("error marshalling index table: %w", err)
	}

	// Creating the byte key.
	key := []byte(indexTableKey)

	// Creating the transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()
	if err := txn.Set(key, itJson); err != nil {
		return fmt.Errorf("error storing index table: %w", err)
	}

	// Committing the transaction to the database.
	if err := txn.Commit(); err != nil {
		return fmt.Errorf("error committing index table: %w", err)
	}

	// If all went all, no errors should be returned.
	return nil
}

// LoadIndexTable is responsible for loading the index table
// from the database.
func (d *Database) LoadIndexTable() (*IndexTable, error) {
	// Converting the key to bytes.
	key := []byte(indexTableKey)

	// Creating the transaction.
	txn := d.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		// If the program is used for the first time, then we won´t have
		// an IndexTable stored in the database. So we need to let that
		// error go and return a new index table to the user, whom will
		// further store the new IndexTable to the database again.
		if errors.Is(err, badger.ErrKeyNotFound) {
			return newIndexTable(), nil
		}

		return nil, fmt.Errorf("error loading index table from database: %w", err)
	}

	// Parsing the item to the expected struct.
	var it IndexTable
	err = item.Value(func(val []byte) error {
		return json.Unmarshal(val, &it)
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing index table: %w", err)
	}

	return &it, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
