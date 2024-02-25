package lights

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/jonamat/makai-automations/pkg/utils"
)

func getEnabled() bool {
	var enabled bool

	err := dbClient.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("light/enabled"))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				setEnabled(DEFAULT_ENABLED)
				item, _ = txn.Get([]byte("light/enabled"))
			} else {
				return err
			}
		}

		err = item.Value(func(val []byte) error {
			enabled = utils.BytesToBool(val)
			return nil
		})
		return err
	})
	if err != nil {
		fmt.Println("Error getting light/enabled: ", err)
	}

	return enabled
}

func setEnabled(value bool) {
	err := dbClient.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("light/enabled"), utils.BoolToBytes(value))
	})
	if err != nil {
		fmt.Println("Error setting light/enabled: ", err)
	}

}

func getLightLevel() int {
	var level int

	err := dbClient.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("light/level"))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				setLightLevel(DEFAULT_LIGHT_LEVEL)
				item, _ = txn.Get([]byte("light/level"))
				return nil
			}
			return err
		}

		err = item.Value(func(val []byte) error {
			level = int(val[0])
			return nil
		})
		return err
	})
	if err != nil {
		fmt.Println("Error getting light/level: ", err)
	}

	return level
}

func setLightLevel(value int) {
	err := dbClient.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("light/level"), []byte{byte(value)})
	})
	if err != nil {
		fmt.Println("Error setting light/level: ", err)
	}
}
