package server

import (
	"github.com/monsterxx03/rkv/config"
	"github.com/go-ini/ini"
	"testing"
	"sync"
	"github.com/go-redis/redis"
	"os"
)

var StartServerOnce sync.Once

const (
	testDataDir string = "/tmp/rocksdb"
	testPort    string = "12345"
)


func GetTestConn() *redis.Client {
	StartServerOnce.Do(RunTestServer)
	return redis.NewClient(&redis.Options{Addr: "127.0.0.1:" + testPort})
}

func RunTestServer() {
	iniCfg := ini.Empty()
	os.RemoveAll(testDataDir)
	iniCfg.Section("server").Key("port").SetValue(testPort)
	iniCfg.Section("rocksdb").Key("data_dir").SetValue(testDataDir)
	cfg := config.NewConfig(iniCfg)
	s := NewServer(cfg)
	go s.Run()
}

func TestServer(t *testing.T) {
	StartServerOnce.Do(RunTestServer)
}
