package alarm

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/jonamat/makai-automations/pkg/utils"
)

const (
	ENABLED_KEY = "alarm-enabled"
)

func getEnabled() (bool, error) {
	var enabled bool

	err := dbClient.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(ENABLED_KEY))
		if err != nil {
			return err

		}

		err = item.Value(func(val []byte) error {
			enabled = utils.BytesToBool(val)
			return nil
		})
		return err
	})

	return enabled, err
}

func setEnabled(value bool) error {
	err := dbClient.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(ENABLED_KEY), utils.BoolToBytes(value))
	})

	return err
}

func setupDb(db *badger.DB) {
	err := db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(ENABLED_KEY))
		if err == badger.ErrKeyNotFound {
			err = txn.Set([]byte(ENABLED_KEY), []byte(utils.BoolToBytes(DEFAULT_ENABLED)))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
