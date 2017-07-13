package main

import (
	"github.com/monsterxx03/rkv/server"
	"github.com/monsterxx03/rkv/config"
	"flag"
	"github.com/go-ini/ini"
	"strconv"
)

var configFile = flag.String("config", "", "path to config file")
var port = flag.Int("port", 0, "port to listen")

func main() {
	flag.Parse()
	var cfg *ini.File
	var err error
	if *configFile != "" {
		// load config from file
		cfg, err = ini.Load(*configFile)
		if err != nil {
			panic(err)
		}
	} else {
		cfg = ini.Empty()
	}
	if *port > 0 {
		cfg.Section("server").Key("port").SetValue(strconv.Itoa(*port))
	}
	s := server.NewServer(config.NewConfig(cfg))
	s.Run()
}
