package db

import (
	_ "github.com/monsterxx03/rkv/db/mysql"
	_ "github.com/monsterxx03/rkv/db/rocksdb"
	"github.com/monsterxx03/rkv/db/backend"
)


type DB struct {
	db backend.IDB
}

func NewDB() *DB {
	b := backend.BackendMap["rocksdb"]
	db, err := b.Open()
	if err != nil {
		panic(err)
	}
	return &DB{db}
}

func (db *DB) Put(key, value []byte) error {
	if err := db.db.Put(key, value); err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	value, err := db.db.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (db *DB) NewBatch() backend.IBatch {
	return db.db.NewBatch()
}

