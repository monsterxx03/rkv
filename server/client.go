package server

import "github.com/monsterxx03/rkv/db"

type client struct {
	cmd string
	args [][]byte
	db *db.DB
	respWriter *RESPWriter
}

func newClient() *client {
	c := new(client)
	return c
}