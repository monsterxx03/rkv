package server

import (
	"net"
	"sync"
	"fmt"
	"log"
	"bufio"
	"github.com/monsterxx03/rkv/db"
	_ "net/http/pprof"
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
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8080", nil))
	//}()
	for {
		select {
		// TODO send quit via signal
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			c := newClient()
			c.server = s
			c.conn = conn
			c.db = s.db
			c.respWriter = NewRESPWriter(conn, s.cfg.WriterBufSize)
			buf := bufio.NewReaderSize(conn, s.cfg.ReaderBufSize)
			c.respReader = NewRESPReader(buf)

			go c.run()
		}
	}
}

func NewServer() *Server {
	// TODO read from cfg file
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

