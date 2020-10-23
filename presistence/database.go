package presistence

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type IDatabase interface {
	AddBlock()
	GetBlock(blockHash string) []byte
	Tip() []byte
	AddTip(blockHash string)
}
type Database struct {
	instance *bolt.DB
	DbFile   string
	bucket   string
}

func NewDatabase(nodeID string) *Database {
	const blocksBucket = "blocks"
	DbFile := fmt.Sprintf("db/blockchain_%s.db", nodeID)
	db, err := bolt.Open(DbFile, 0600, &bolt.Options{Timeout: 2 * time.Minute})
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			_, _ = tx.CreateBucket([]byte(blocksBucket))
		}
		return nil
	})
	newDatabase := &Database{
		db,
		DbFile,
		blocksBucket,
	}
	return newDatabase
}

func (db *Database) AddBlock(serializedBlock []byte, blockHash string) {
	_ = db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte(blockHash), serializedBlock)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
}
func (db *Database) GetBlock(blockHash string) []byte {
	var serializedBlock []byte
	err := db.instance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		serializedBlock = b.Get([]byte(blockHash))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return serializedBlock
}

func (db *Database) Tip() []byte {
	var serializedBlock []byte
	err := db.instance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		serializedBlock = b.Get([]byte("t"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return serializedBlock
}
func (db *Database) AddTip(blockHash string) {
	_ = db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte("t"), []byte(blockHash))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
}
