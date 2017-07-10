package backend

import (
	"fmt"
)

type IBackend interface {
	String() string
	Open() (IDB, error)
}

type IDB interface{
	Close() error
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) (error)
}

var BackendMap = map[string]IBackend{}

func RegisterBackend(b IBackend) {
	name := b.String()
	if _, ok := BackendMap[name]; ok {
		panic(fmt.Errorf("Backend '%s' has been registered", name))
	}
	BackendMap[name] = b
}
