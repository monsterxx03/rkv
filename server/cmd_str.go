package server

import (
	"github.com/monsterxx03/rkv/codec"
)

func cmdGet(c *client) error {
	if len(c.args) != 1 {
		return &WrongParamError{"get"}
	}
	value, err := c.db.Get(c.args[0])
	if err != nil {
		return err
	}
	c.respWriter.writeStr(string(value)) // TODO change to real write resp
	return nil
}

func cmdSet(c *client) error {
	if len(c.args) != 2 {
		return &WrongParamError{"set"}
	}
	if err := c.db.Put(c.args[0], codec.EncodeStrVal(c.args[1])); err != nil {
		return err
	}
	return nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
}
