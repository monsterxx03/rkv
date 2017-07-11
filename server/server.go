package server

import (
	"net"
	"sync"
	"fmt"
	"log"
	"runtime"
	"bufio"
	"errors"
	"strings"
	"github.com/monsterxx03/rkv/db"
)

const (
	DefaultAddr          string = "0.0.0.0"
	DefaultPort          int    = 9910
	DefaultReaderBufSize        = 4096
	DefaultWriterBufSize        = 4096
)

type Config struct {
	Addr          string
	Port          int
	ReaderBufSize int
	WriterBufSize int
}

type Server struct {
	cfg      *Config
	listener net.Listener
	db       *db.DB
	wg       sync.WaitGroup
	quit     <-chan struct{}
}

func (s *Server) Run() {
	for {
		select {
		// TODO send quit via signal
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				panic(err)
			}
			go handleReq(conn, s)
		}
	}
}

func NewServer() *Server {
	cfg := &Config{DefaultAddr, DefaultPort, DefaultReaderBufSize, DefaultWriterBufSize}
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	log.Println("Listening at:", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Failed to listen on: " + addr)
	}
	return &Server{
		cfg: cfg, listener: listener,
		db:  db.NewDB(), quit: make(chan struct{})}
}

func handleReq(conn net.Conn, serv *Server) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			length := runtime.Stack(buf, false)
			buf = buf[0:length]
			log.Fatalf("panic when handleReq: %s, %v", buf, err)
		}
		conn.Close()
		serv.wg.Done()
	}()
	serv.wg.Add(1)
	buf := bufio.NewReaderSize(conn, serv.cfg.ReaderBufSize)
	reader := NewRESPReader(buf)
	for {
		// continue read from client
		data, err := reader.ParseRequest()
		if err != nil {
			log.Println(err)
			return
		}
		c := newClient()
		c.db = serv.db
		c.cmd = strings.ToLower(string(data[0]))
		c.args = data[1:]
		c.respWriter = NewRESPWriter(conn, serv.cfg.WriterBufSize)
		err = handleCmd(c)
		if err != nil {
			c.respWriter.writeError(err)
		}
		c.respWriter.flush()
	}
}

func handleCmd(c *client) error {
	cmdStr := c.cmd
	cmdFunc, ok := CommandsMap[cmdStr]
	if !ok {
		return errors.New("Err unknown command " + cmdStr)
	}
	return cmdFunc(c)
}
