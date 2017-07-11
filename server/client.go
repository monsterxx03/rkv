package server

import (
	"github.com/monsterxx03/rkv/db"
	"sync"
)

type client struct {
	cmd string
	args [][]byte
	db *db.DB
	respWriter *RESPWriter
}

type KeyGuard struct {
	key []byte
	m sync.Mutex
}

func newClient() *client {
	c := new(client)
	return c
}