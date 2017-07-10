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
)

const (
	DefaultAddr          string = "0.0.0.0"
	DefaultPort          int    = 9910
	DefaultReaderBufSize        = 4096
)

type Config struct {
	Addr          string
	Port          int
	ReaderBufSize int
}

type Server struct {
	cfg      *Config
	listener net.Listener
	Wg       sync.WaitGroup
	quit     <-chan struct{}
}

func (s *Server) Run() {
	for {
		select {
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
	cfg := &Config{DefaultAddr, DefaultPort, DefaultReaderBufSize}
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	log.Println("Listening at:", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Failed to listen on: " + addr)
	}
	return &Server{cfg: cfg, listener: listener, quit: make(chan struct{})}
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
		serv.Wg.Done()
	}()
	serv.Wg.Add(1)
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
		c.cmd = strings.ToLower(string(data[0]))
		c.args = data[1:]
		result, err := handleCmd(c)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err.Error())))
		} else {
			conn.Write(result) // fake response
		}
	}
}

func handleCmd(c *client) ([]byte, error) {
	cmdStr := c.cmd
	cmdFunc, ok := CommandsMap[cmdStr]
	if !ok {
		return nil, errors.New("unknown command " + cmdStr)
	}
	function := cmdFunc
	response, err := function(c)
	if err != nil {
		return nil, err
	}
	return response, nil
}
