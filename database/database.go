package database

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
type DB struct {
	instance *bolt.DB
	DbFile   string
	bucket   string
}

func NewDatabase(nodeID string) *DB {
	const blocksBucket = "blocks"
	DbFile := fmt.Sprintf("database_%s.db", nodeID)
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
	newDatabase := &DB{
		db,
		DbFile,
		blocksBucket,
	}
	return newDatabase
}

func (db *DB) AddBlock(serializedBlock []byte, blockHash string) {
	_ = db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte(blockHash), serializedBlock)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
}
func (db *DB) GetBlock(blockHash string) []byte {
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

func (db *DB) Tip() []byte {
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
func (db *DB) AddTip(blockHash string) {
	_ = db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.bucket))
		err := b.Put([]byte("t"), []byte(blockHash))
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
}
