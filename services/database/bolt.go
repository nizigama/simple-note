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

func Store(item map[string]interface{}, tableName string) error {

	var id []byte = make([]byte, 8)
	xi, err := json.Marshal(item)

	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(tableName))

		if b == nil {
			return fmt.Errorf("No such table found")
		}

		nxtID, err := b.NextSequence()

		if err != nil {
			return err
		}

		binary.LittleEndian.PutUint32(id, uint32(nxtID))

		return b.Put(id, xi)
	})
}

// CloseDBConnection closes the database connection and releases the db so that other apps can use it
func CloseDBConnection() {
	db.Close()
}
