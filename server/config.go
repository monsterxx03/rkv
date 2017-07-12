package server


type Config struct {
	Addr          string
	Port          int
	ReaderBufSize int
	WriterBufSize int
}


func newConfig() *Config  {
	return &Config{DefaultAddr, DefaultPort, DefaultReaderBufSize, DefaultWriterBufSize}
}

