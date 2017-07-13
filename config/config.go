package config

import "github.com/go-ini/ini"

type Config struct {
	Addr         string
	Port         int
	ReadBufSize  int
	WriteBufSize int
	RocksDB      RocksDBConfig
}

type RocksDBConfig struct {
	DataDir                       string
	BlockCache                    int
	BlockSize                     int
	BloomFilterBitsPerKey         int
	BackgroundThreads             int
	HighPriorityBackgroundThreads int
	CompressionType               uint
	WriteBufferSize               int
	MaxWriteBufferNumber          int
	MinWriteBufferNumberToMerge   int
	MaxOpenFiles                  int
	NumLevels                     int
	MaxBackgroundCompactions      int
	MaxBackgroundFlushes          int
	UseFsync                      bool
}

func NewConfig(cfg *ini.File) *Config {
	_cfg := new(Config)
	_cfg.Addr = cfg.Section("server").Key("host").MustString("127.0.0.1")
	_cfg.Port = cfg.Section("server").Key("port").MustInt(12000)
	_cfg.ReadBufSize = cfg.Section("server").Key("read_buf_size").MustInt(4096)
	_cfg.WriteBufSize = cfg.Section("server").Key("write_buf_size").MustInt(4096)
	_cfg.RocksDB = newRocksDBConfig(cfg.Section("rocksdb"))
	return _cfg
}

func newRocksDBConfig(section *ini.Section) RocksDBConfig {
	return RocksDBConfig{
		DataDir:                       section.Key("data_dir").MustString("data"),
		BlockCache:                    section.Key("block_cache").MustInt(1073741824),
		BlockSize:                     section.Key("block_size").MustInt(65536),
		BloomFilterBitsPerKey:         section.Key("bloomfilter_bitsperkey").MustInt(10),
		BackgroundThreads:             section.Key("background_threads").MustInt(16),
		HighPriorityBackgroundThreads: section.Key("high_priority_background_threads").MustInt(1),
		CompressionType:               section.Key("compression_type").MustUint(0),
		WriteBufferSize:               section.Key("write_buffer_size").MustInt(134217728),
		MaxWriteBufferNumber:          section.Key("max_write_buffer_number").MustInt(6),
		MinWriteBufferNumberToMerge:   section.Key("min_write_buffer_number_to_merge").MustInt(2),
		MaxOpenFiles:                  section.Key("max_open_files").MustInt(1024),
		NumLevels:                     section.Key("num_levels").MustInt(7),
		MaxBackgroundCompactions:      section.Key("max_background_compactions").MustInt(15),
		MaxBackgroundFlushes:          section.Key("max_background_flushes").MustInt(1),
		UseFsync:                      section.Key("use_fsync").MustBool(false),
	}
}
