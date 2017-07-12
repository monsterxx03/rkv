package server

import (
	"github.com/monsterxx03/rkv/codec"
	"strconv"
	"log"
	"errors"
)

func cmdGet(c *client, args Args) error {
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

func cmdSet(c *client, args Args) error {
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

func cmdIncr(c *client, args Args) error {
	if len(args) != 1 {
		return &WrongParamError{"incr"}
	}
	c.db.Locker.Lock(string(args[0]))
	defer c.db.Locker.Unlock(string(args[0]))

	if value, err := c.db.Get(args[0]); err != nil {
		return err
	} else if len(value) > 0 && codec.DecodeType(value) != codec.StrType {
		return WrongTypeError
	} else {
		value := codec.DecodeStrKey(value)
		// try to convert value to int
		n, err := strconv.ParseInt(string(value), 10, 64)
		if err != nil {
			log.Println(err)
			return errors.New("ERR value is not an integer or out of range")
		}
		n += 1
		if err := c.db.Put(args[0], codec.EncodeStrVal(Int64ToSlice(n))) ; err != nil {
			return err
		}
		c.respWriter.writeInt(n)
	}
	return nil
}

func cmdDecr(c *client, args Args) error {
	return nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
	register("incr", cmdIncr)
	register("decr", cmdDecr)
}
