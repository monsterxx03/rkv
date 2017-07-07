package redis

var CommandsMap = map[string]interface {} {
	"info": cmdInfo,
	"ping": cmdPing,
	"echo": cmdEcho,
}

func cmdInfo(data [][]byte) ([]byte, error) {
	return []byte("+INFO\r\n"), nil
}

func cmdPing(data [][]byte) ([]byte, error) {
	return []byte("+PONG\r\n"), nil
}

func cmdEcho(data [][]byte) ([]byte, error) {
	return []byte("+hahah\r\n"), nil
}

