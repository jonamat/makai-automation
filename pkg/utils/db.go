package utils

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

var dbPath = "data/data.bdg"

func CreateDbClient() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db
}
