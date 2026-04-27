package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	badger "github.com/dgraph-io/badger/v4"
)

// Database encapsulates the badger.DB type.
type Database struct {
	db *badger.DB
}

// NewDatabase is the Database constructor. It takes in
// the directory path used to store the data and returns
// a pointer to Database and may return an error.
func NewDatabase(dirPath string) (*Database, error) {
	db, err := badger.Open(badger.DefaultOptions(dirPath))
	if err != nil {
		return nil, fmt.Errorf("unable to open badgerDB: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

// storeObject will be used to store an object to the database.
func (d *Database) storeObject(key string, value interface{}) error {
	var prefix string

	// Getting the object type.
	t := reflect.TypeOf(value)
	switch t.Name() {
	case "Manifest":
		prefix = "mani:"
	case "IndexTable":
		prefix = "it:"
	}

	// Starting the transaction.
	txn := d.db.NewTransaction(true)
	defer txn.Discard()

	// Marshalling the value.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	// Making the transaction
	dbKey := prefix + key
	if err := txn.Set([]byte(dbKey), buf.Bytes()); err != nil {
		return fmt.Errorf("error making transaction: %w", err)
	}

	// Commiting the transaction.
	if err := txn.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	return nil
}
