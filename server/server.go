package server

import (
	"net"
	"sync"
	"fmt"
	"log"
	"bufio"
	"github.com/monsterxx03/rkv/db"
	"github.com/monsterxx03/rkv/config"
	_ "net/http/pprof"
)

type Server struct {
	cfg      *config.Config
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
			c.respWriter = NewRESPWriter(conn, s.cfg.WriteBufSize)
			buf := bufio.NewReaderSize(conn, s.cfg.ReadBufSize)
			c.respReader = NewRESPReader(buf)

			go c.run()
		}
	}
}

func NewServer(cfg *config.Config) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	log.Println("Listening at:", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("Failed to listen on: " + addr)
	}
	_db := db.NewDB(cfg)
	return &Server{
		cfg: cfg, listener: listener,
		db:  _db, quit: make(chan struct{})}
}

