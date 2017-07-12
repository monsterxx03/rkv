package server

import (
	"github.com/monsterxx03/rkv/db"
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
	cmd string
	db *db.DB
	respReader *RESPReader
	respWriter *RESPWriter
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
		c.cmd = cmd
		if err := c.exeCmd(data[1:]) ; err != nil {
			c.respWriter.writeError(err)
		}
		if err := c.respWriter.flush(); err != nil {
			log.Println("Error flushing response: ", err)
		}
	}
}

func (c *client) exeCmd(args [][]byte) error {
	cmdFunc, ok := CommandsMap[c.cmd]
	if !ok {
		return errors.New("Err unknown command " + c.cmd)
	}
	if err := cmdFunc(c, args); err != nil {
		return err
	}
	return nil
}

