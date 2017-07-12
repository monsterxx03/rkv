package db

import (
	_ "github.com/monsterxx03/rkv/db/mysql"
	_ "github.com/monsterxx03/rkv/db/rocksdb"
	"github.com/monsterxx03/rkv/db/backend"
	"sync"
)


type DB struct {
	db backend.IDB
	WriteBatch backend.IBatch
	Locker backend.ILock
}

func NewDB() *DB {
	b := backend.BackendMap["rocksdb"]
	db, err := b.Open()
	if err != nil {
		panic(err)
	}
	return &DB{db, db.NewBatch(), newLock()}
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


// TODO clean up locks to avoid memory leak
type MemLock struct {
	locks map[string]*sync.Mutex
}

func (l *MemLock) Lock(key string) {
	if _l, ok := l.locks[key]; !ok {
		_ll := new(sync.Mutex)
		_ll.Lock()
		l.locks[key] = _ll
	} else {
		_l.Lock()
	}
}

func (l *MemLock) Unlock(key string) {
	if _l, ok := l.locks[key]; ok {
		_l.Unlock()
	}
}

func newLock() backend.ILock {
	// Only support lock in memory
	l := new(MemLock)
	l.locks = make(map[string]*sync.Mutex)
	return l
}
