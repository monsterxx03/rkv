package mysql

import (
	"github.com/monsterxx03/rkv/db/backend"
)

type Backend struct {
}

func (s Backend) String() string {
	return "mysql"
}

func (s Backend) Open() (backend.IDB, error) {
	db := new(DB)
	if err := db.open(); err != nil {
		return nil, err
	}
	return db, nil
}

type DB struct {
}

func (db *DB) open() error {
	return nil
}

func (db *DB) Close() error {
	return nil
}

func (db *DB) Put(key, value []byte) error {
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return nil, nil
}

func (db *DB) Delete(key []byte) error {
	return nil
}

func init() {
	backend.RegisterBackend(Backend{})
}
