package utils

import (
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

var dbPath = "data/"
var dbClient *badger.DB

func CreateDbClient() *badger.DB {
	if dbClient != nil {
		return dbClient
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatal(err)
	}

	dbClient = db

	return db
}
