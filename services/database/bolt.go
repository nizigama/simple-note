package boltDB

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	db  *bolt.DB = nil
	err error
)

// InitDBConnection initiates the connection with a database
func InitDBConnection(logger *log.Logger, tables ...string) {
	db, err = bolt.Open("simpnote.db", 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		logger.Fatalln("Initiate database connection:", err)
	}

	// checking if there are tables, and if there aren't create them
	err = db.Update(func(tx *bolt.Tx) error {

		for _, table := range tables {
			if _, tErr := tx.CreateBucketIfNotExists([]byte(table)); tErr != nil {
				return tErr
			}
		}

		return nil
	})

	if err != nil {
		logger.Fatalln("Create database tables:", err)
	}
}

func Store(item map[string]interface{}, tableName string) (uint64, error) {

	var id []byte = make([]byte, 8)
	xi, err := json.Marshal(item)

	if err != nil {
		return 0, err
	}

	var userID uint64

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		userID, err = b.NextSequence()

		if err != nil {
			return err
		}

		binary.LittleEndian.PutUint32(id, uint32(userID))

		return b.Put(id, xi)
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func Show(itemID uint64, tableName string) (map[string]interface{}, error) {

	var itemData map[string]interface{}
	result := []byte{}

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		id := make([]byte, 8)

		binary.LittleEndian.PutUint32(id, uint32(itemID))

		result = b.Get(id)
		if result == nil {
			return fmt.Errorf("no such item found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result, &itemData)

	if err != nil {
		return nil, err
	}

	return itemData, nil
}

func All(tableName string) ([]map[string]interface{}, error) {

	var itemsData []map[string]interface{}

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		err = b.ForEach(func(k []byte, v []byte) error {

			itemData := map[string]interface{}{}

			err = json.Unmarshal(v, &itemData)

			if err != nil {
				return err
			}
			itemData["itemID"] = binary.LittleEndian.Uint64(k)

			itemsData = append(itemsData, itemData)

			return nil
		})
		if err != nil {
			return fmt.Errorf("no such item found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return itemsData, nil
}

func SingleByStringField(tableName, fieldName, fieldValue string) (map[string]interface{}, uint64, error) {

	var matchingItemData map[string]interface{}
	var key uint64

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		err = b.ForEach(func(k []byte, v []byte) error {

			itemData := map[string]interface{}{}

			err = json.Unmarshal(v, &itemData)

			if err != nil {
				return err
			}

			if itemData[fieldName] == fieldValue {
				key = binary.LittleEndian.Uint64(k)
				matchingItemData = itemData
			}

			return nil
		})
		if err != nil || key == 0 {
			return fmt.Errorf("no such item found")
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return matchingItemData, key, nil
}

func ManyByStringField(tableName, fieldName, fieldValue string) ([]map[string]interface{}, error) {

	var matchingItemsData []map[string]interface{}

	err := db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		err = b.ForEach(func(k []byte, v []byte) error {

			itemData := map[string]interface{}{}

			err = json.Unmarshal(v, &itemData)

			if err != nil {
				return err
			}

			if itemData[fieldName] == fieldValue {
				matchingItemsData = append(matchingItemsData, itemData)
			}

			return nil
		})
		if err != nil || len(matchingItemsData) == 0 {
			return fmt.Errorf("no such item found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchingItemsData, nil
}

func Update(item map[string]interface{}, tableName string, itemID uint64) error {

	xi, err := json.Marshal(item)

	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		id := make([]byte, 8)

		binary.LittleEndian.PutUint32(id, uint32(itemID))

		result := b.Get(id)

		if result == nil {
			return fmt.Errorf("no item found")
		}

		err = b.Delete(id)

		if err != nil {
			return fmt.Errorf("error deleting old item")
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		id := make([]byte, 8)

		binary.LittleEndian.PutUint32(id, uint32(itemID))

		return b.Put(id, xi)
	})

	return nil
}

func Delete(tableName string, itemID uint64) error {

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("no such table found")
		}

		id := make([]byte, 8)

		binary.LittleEndian.PutUint32(id, uint32(itemID))

		result := b.Get(id)

		if result == nil {
			return fmt.Errorf("no item found")
		}

		err = b.Delete(id)

		if err != nil {
			return fmt.Errorf("error deleting item")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// CloseDBConnection closes the database connection and releases the db so that other apps can use it
func CloseDBConnection() {
	db.Close()
}
