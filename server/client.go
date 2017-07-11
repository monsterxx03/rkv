package server

import (
	"github.com/monsterxx03/rkv/db"
	"sync"
	"github.com/monsterxx03/rkv/db/backend"
	"io"
	"log"
	"strings"
	"runtime"
	"net"
)

type client struct {
	server *Server
	conn net.Conn
	cmd string
	args [][]byte
	db *db.DB
	respReader *RESPReader
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

func (c *client) run() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			length := runtime.Stack(buf, false)
			buf = buf[0:length]
			log.Printf("panic when handleReq: %s, %v", buf, err)
		}
		c.conn.Close()
		c.server.wg.Done()
	}()
	c.server.wg.Add(1)
	for {
		data, err := c.respReader.ParseRequest()
		if err != nil {
			if err == io.EOF {
				// client close connection
				return
			} else {
				// unexpected error
				panic(err)
			}
		}
		c.cmd = strings.ToLower(string(data[0]))
		c.args = data[1:]
		if err := handleCmd(c); err != nil {
			c.respWriter.writeError(err)
		}
		c.respWriter.flush()
	}
}