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

	if len(value) > 0 && codec.DecodeType(value) != codec.StrType {
			return WrongTypeError
	}

	value = codec.DecodeStrKey(value)
	c.respWriter.writeBulkStr(value)
	return nil
}

func cmdSet(c *client) error {
	if len(c.args) != 2 {
		return &WrongParamError{"set"}
	}
	/*
	if value, err := c.db.Get(c.args[0]); err != nil {
		return err
	} else if len(value) > 0 && codec.DecodeType(value) != codec.StrType {
			return WrongTypeError
	}
	*/
	if err := c.db.Put(c.args[0], codec.EncodeStrVal(c.args[1])); err != nil {
		return err
	}
	c.respWriter.writeStr("OK")
	return nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
}
