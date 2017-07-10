package main

import (
	"github.com/monsterxx03/rkv/server"
)

func main() {
	s := server.NewServer()
	s.Run()
}
