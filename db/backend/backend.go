package backend

import (
	"fmt"
	"github.com/monsterxx03/rkv/config"
)

type IBackend interface {
	String() string
	Open(*config.Config) (IDB, error)
}

type IBatch interface {
	Close()
	Put(key, value []byte)
	Delete(key []byte)
	Commit() error
}

type IIter interface {
	Next()
}

type IDB interface {
	Close() error
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	MGet(keys [][]byte) ([][]byte, error)
	Delete(key []byte) (error)
	NewBatch() IBatch
	NewIter() IIter
}

type ILock interface {
	Lock(key string)
	Unlock(key string)
}

var BackendMap = map[string]IBackend{}

func RegisterBackend(b IBackend) {
	name := b.String()
	if _, ok := BackendMap[name]; ok {
		panic(fmt.Errorf("Backend '%s' has been registered", name))
	}
	BackendMap[name] = b
}

