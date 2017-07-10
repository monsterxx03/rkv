package server

import (
	"github.com/monsterxx03/rkv/codec"
)

func cmdGet(c *client) ([]byte, error) {
	if len(c.args) != 1 {
		return nil, &WrongParamError{"get"}
	}
	value, err := c.db.Get(c.args[0])
	if err != nil {
		return nil, err
	}
	return value, nil
}

func cmdSet(c *client) ([]byte, error) {
	if len(c.args) != 2 {
		return nil, &WrongParamError{"set"}
	}
	if err := c.db.Put(c.args[0], codec.EncodeStrVal(c.args[1])); err != nil {
		return nil, err
	}
	return []byte("+OK\r\n"), nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
}
