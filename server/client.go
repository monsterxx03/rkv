package server

import (
	"github.com/monsterxx03/rkv/db"
	"sync"
	"github.com/monsterxx03/rkv/db/backend"
)

type client struct {
	cmd string
	args [][]byte
	db *db.DB
	respWriter *RESPWriter
	writeBatch backend.IBatch
}

type KeyGuard struct {
	key []byte
	m sync.Mutex
}

func newClient() *client {
	c := new(client)
	return c
}