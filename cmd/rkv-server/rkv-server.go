package main

import (
	"github.com/monsterxx03/rkv/server"
	"flag"
)

var configFile = flag.String("config", "", "path to config file")
var port = flag.Int("port", 12000, "port to listen")

func main() {
	flag.Parse()
	if *configFile != "" {
		// load config from file
	}
	// override config from cmd line
	s := server.NewServer()
	s.Run()
}
