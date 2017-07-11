package server

import "fmt"


type CommandFunc func(*client) error

var CommandsMap = map[string]CommandFunc{}

func register(cmd string, function CommandFunc) {
	_, ok := CommandsMap[cmd]
	if ok {
		panic(fmt.Errorf("cmd '%s' has been registered", cmd))
	}
	CommandsMap[cmd] = function
}

func cmdInfo(c *client) error {
	c.respWriter.writeStr("info")
	return nil
}

func cmdPing(c *client) error {
	c.respWriter.writeStr("pong")
	return nil
}

func cmdEcho(c *client) error {
	c.respWriter.writeStr("hahah")
	return nil
}

func init() {
	register("info", cmdInfo)
	register("ping", cmdPing)
	register("echo", cmdEcho)
}
