package server

import (
	"github.com/monsterxx03/rkv/db"
	"sync"
	"io"
	"log"
	"strings"
	"runtime"
	"net"
	"errors"
)

type client struct {
	server *Server
	conn net.Conn
	db *db.DB
	respReader *RESPReader
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
		cmd := strings.ToLower(string(data[0]))
		cmdFunc, ok := CommandsMap[cmd]
		if !ok {
			c.respWriter.writeError(errors.New("Err unknown command " + cmd))
		}
		if err != c.exeCmd(cmdFunc, data[1:]) {
			log.Println("Error flushing response: ", err)
		}
	}
}

func (c *client) exeCmd(cmd CommandFunc, args [][]byte) error {
	if err := cmd(c, args); err != nil {
		c.respWriter.writeError(err)
	}
	if err := c.respWriter.flush(); err != nil {
		return err
	}
	return nil
}

