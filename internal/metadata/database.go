package metadata

import (
	"fmt"

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

// func (d *Database) storeObj(key interface{}) (error) {
// 	err := d.db.Update(func(txn *badger.Txn) error {
		
// 		return nil
// 	})
// 	return nil
// }


