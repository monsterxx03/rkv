package db

import (
	_ "github.com/monsterxx03/rkv/db/mysql"
	_ "github.com/monsterxx03/rkv/db/rocksdb"
	"github.com/monsterxx03/rkv/db/backend"
)


type DB struct {
	db backend.IDB
	writeBatch backend.IBatch
}

func NewDB() *DB {
	b := backend.BackendMap["rocksdb"]
	db, err := b.Open()
	if err != nil {
		panic(err)
	}
	return &DB{db, db.NewBatch()}
}

func (db *DB) Put(key, value []byte) error {
	if err := db.db.Put(key, value); err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	if value, err := db.db.Get(key); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}
