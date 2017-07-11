package server

import (
	"fmt"
	"github.com/monsterxx03/rkv/codec"
)

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

func cmdType(c *client) error {
	if len(c.args) != 1 {
		return &WrongParamError{"type"}
	}
	value, err := c.db.Get(c.args[0])
	if err != nil {
		return err
	}
	if len(value) == 0 {
		c.respWriter.writeStr("none")
		return nil
	}
	switch dataType := codec.DecodeType(value); dataType {
	case codec.StrType:
		c.respWriter.writeStr("string")
	case codec.ListType:
		c.respWriter.writeStr("list")
	case codec.HashType:
		c.respWriter.writeStr("hash")
	case codec.SetType:
		c.respWriter.writeStr("set")
	case codec.ZSetType:
		c.respWriter.writeStr("zset")
	default:
		c.respWriter.writeStr("unknown type " + string(dataType))
	}
	return nil
}

// Get key first, delete based on its data type
func cmdDel(c *client) error {
	return nil
}

func init() {
	register("info", cmdInfo)
	register("ping", cmdPing)
	register("echo", cmdEcho)
	register("del", cmdDel)
	register("type", cmdType)
}
