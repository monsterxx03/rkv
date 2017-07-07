package main

import (
	"fmt"
	"net"
	"log"
	"bufio"
	"github.com/monsterxx03/rkv/redis"
	"sync"
	"runtime"
	"errors"
)

const (
	DefaultAddr string = "0.0.0.0"
	DefaultPort int = 9910
	DefaultReaderBufSize = 4096
)

type Config struct {
	Addr string
	Port int
	ReaderBufSize int
}

type Server struct {
	cfg *Config
	listener net.Listener
	Wg sync.WaitGroup
	quit <- chan struct{}
}

func (s *Server) Run() {
	for {
		select {
		case <- s.quit:
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


func handleReq(conn net.Conn, server *Server) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			length := runtime.Stack(buf, false)
			buf = buf[0:length]
			log.Fatalf("panic when handleReq: %s, %v", buf, err)
		}
		conn.Close()
		server.Wg.Done()
	}()
	server.Wg.Add(1)
	buf := bufio.NewReaderSize(conn, server.cfg.ReaderBufSize)
	reader := redis.NewRESPReader(buf)
	for {
		// continue read from client
		data, err := reader.ParseRequest()
		if err != nil {
			log.Println(err)
			return
		}
		// TODO handle error
		result, _ := handleCmd(data[0], data[1:])
		conn.Write(result)  // fake response
	}
}

func handleCmd(cmd []byte, args [][]byte) ([]byte, error) {
	cmdStr := string(cmd)
	cmdFunc, ok := redis.CommandsMap[cmdStr]
	if !ok {
		return nil, errors.New("unknown command" + cmdStr)
	}
	function := cmdFunc.(func ([][]byte) ([]byte, error))
	response, err := function(args)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func main() {
	server := NewServer()
	server.Run()
}