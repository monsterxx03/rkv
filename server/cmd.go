package server

import (
	"fmt"
	"github.com/monsterxx03/rkv/codec"
)

type Args [][]byte

// return key name's string format
func (a Args) skey() string {
	return string(a[0])
}

func (a Args) key() []byte {
	return a[0]
}

// return first value
func (a Args) value() []byte {
	return a[1]
}

// return all values except first one
func (a Args) values() [][]byte {
	return a[1:]
}

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
	value, err := c.db.Get(args.key())
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
	if len(args) == 0 {
		return &WrongParamError{"del"}
	}

	values, err := c.db.MGet(args)
	if err != nil {
		return err
	}
	// filter empty value
	// values = FilterByte(values, func (x []byte) bool { return len(x) > 0})

	if AllByte(values, func (x []byte) bool {return len(x) == 0}) {
		// all keys not exist
		c.respWriter.writeInt(0)
		return nil
	}
	batch := c.db.NewBatch()
	for i, key := range args {
		c.Lock(string(key))
		defer c.Unlock(string(key))

		if len(values[i]) == 0 {
			// skip non-exist key
			continue
		}
		switch dataType := codec.DecodeType(values[i]); dataType {
		case codec.StrType:
			if err := delStr(batch, key, values[i]); err != nil {
				return err
			}
		case codec.ListType:
			if err := delList(batch, key, values[i]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown data type %s", dataType)
		}
	}
	batch.Commit()
	c.respWriter.writeInt(1)
	return nil
}

func init() {
	register("info", cmdInfo)
	register("ping", cmdPing)
	register("echo", cmdEcho)
	register("del", cmdDel)
	register("type", cmdType)
}
