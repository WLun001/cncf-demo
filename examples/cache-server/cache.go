package main

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"time"
)

type CacheItem struct {
	Body      []byte `json:"body"`
	Ctype     []byte `json:"ctype"`
	Cencoding []byte `json:"cencoding"`
	Status    int    `json:"status"`
}

type CacheStorage struct {
	db *badger.DB
}

func NewCacheStorage() CacheStorage {
	// in memory mode
	// https://dgraph.io/docs/badger/get-started/#in-memory-mode-diskless-mode
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		log.Fatal(err)
	}
	return CacheStorage{db: db}
}

// Get value by key
func (c *CacheStorage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	var data []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		// item.Value() is only valid within the transaction.
		// We can either copy it ourselves or use the ValueCopy() method.
		// or to keep the tx open until unmarshalling is done.
		data, err = item.ValueCopy(nil)
		return err
	})
	// If no value was found return false
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return data, err
}

// Set key with value
func (c *CacheStorage) Set(key string, val []byte, exp time.Duration) error {
	// Ain't Nobody Got Time For That
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	entry := badger.NewEntry(utils.UnsafeBytes(key), val)
	if exp != 0 {
		entry.WithTTL(exp)
	}
	return c.db.Update(func(tx *badger.Txn) error {
		return tx.SetEntry(entry)
	})
}
