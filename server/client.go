package server

import "github.com/monsterxx03/rkv/db"

type client struct {
	cmd string
	args [][]byte
	db *db.DB
}

func newClient() *client {
	c := new(client)
	c.db = db.NewDB()
	return c
}