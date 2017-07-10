package server

import "fmt"


type CommandFunc func(*client) ([]byte, error)

var CommandsMap = map[string]CommandFunc{}

func register(cmd string, function CommandFunc) {
	_, ok := CommandsMap[cmd]
	if ok {
		panic(fmt.Errorf("cmd '%s' has been registered", cmd))
	}
	CommandsMap[cmd] = function
}

func cmdInfo(c *client) ([]byte, error) {
	return []byte("+INFO\r\n"), nil
}

func cmdPing(c *client) ([]byte, error) {
	return []byte("+PONG\r\n"), nil
}

func cmdEcho(c *client) ([]byte, error) {
	return []byte("+hahah\r\n"), nil
}

func init() {
	register("info", cmdInfo)
	register("ping", cmdPing)
	register("echo", cmdEcho)
}
