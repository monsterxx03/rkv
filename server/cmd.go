package server

import (
	"fmt"
	"github.com/monsterxx03/rkv/codec"
)

type Args [][]byte
type CommandFunc func(*client, Args) error

var CommandsMap = map[string]CommandFunc{}

func register(cmd string, function CommandFunc) {
	_, ok := CommandsMap[cmd]
	if ok {
		panic(fmt.Errorf("cmd '%s' has been registered", cmd))
	}
	CommandsMap[cmd] = function
}

func cmdInfo(c *client, args Args) error {
	c.respWriter.writeStr("info")
	return nil
}

func cmdPing(c *client, args Args) error {
	c.respWriter.writeStr("pong")
	return nil
}

func cmdEcho(c *client, args Args) error {
	c.respWriter.writeStr("hahah")
	return nil
}

func cmdType(c *client, args Args) error {
	if len(args) != 1 {
		return &WrongParamError{"type"}
	}
	value, err := c.db.Get(args[0])
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
func cmdDel(c *client, args Args) error {
	return nil
}

func init() {
	register("info", cmdInfo)
	register("ping", cmdPing)
	register("echo", cmdEcho)
	register("del", cmdDel)
	register("type", cmdType)
}
