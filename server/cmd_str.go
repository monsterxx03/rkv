package server

import (
	"github.com/monsterxx03/rkv/codec"
)

func cmdGet(c *client, args [][]byte) error {
	if len(args) != 1 {
		return &WrongParamError{"get"}
	}
	value, err := c.db.Get(args[0])
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

func cmdSet(c *client, args [][]byte) error {
	if len(args) != 2 {
		return &WrongParamError{"set"}
	}
	// check key type
	if value, err := c.db.Get(args[0]); err != nil {
		return err
	} else if len(value) > 0 && codec.DecodeType(value) != codec.StrType {
			return WrongTypeError
	}

	if err := c.db.Put(args[0], codec.EncodeStrVal(args[1])); err != nil {
		return err
	}
	c.respWriter.writeStr("OK")
	return nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
}
