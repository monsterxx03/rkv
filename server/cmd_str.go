package server

import (
	"github.com/monsterxx03/rkv/codec"
	"log"
)

func cmdGet(c *client, args Args) error {
	if !args.matchN(1) {
		return &WrongParamError{"get"}
	}
	value, err := c.db.Get(args[0])
	if err != nil {
		return err
	}
	if err := codec.CheckStrType(value); err != nil {
		return err
	}
	value = codec.DecodeStrKey(value)
	c.respWriter.writeBulkStr(value)
	return nil
}

func cmdSet(c *client, args Args) error {
	if !args.matchN(2) {
		return &WrongParamError{"set"}
	}
	c.Lock(args.skey())
	defer c.Unlock(args.skey())

	// check key type
	value, err := c.db.Get(args[0])
	if err != nil {
		return err
	}
	if err := codec.CheckStrType(value); err != nil {
		return err
	}

	if err := c.db.Put(args[0], codec.EncodeStrVal(args[1])); err != nil {
		return err
	}
	c.respWriter.writeStr("OK")
	return nil
}

func cmdIncr(c *client, args Args) error {
	if !args.matchN(1) {
		return &WrongParamError{"incr"}
	}
	c.Lock(args.skey())
	defer c.Unlock(args.skey())

	value, err := c.db.Get(args[0])
	if err != nil {
		return err
	}
	if err := codec.CheckStrType(value); err != nil {
		return err
	}
	value = codec.DecodeStrKey(value)
	// try to convert value to int
	n, err := SliceToInt64(value)
	if err != nil {
		log.Println(err)
		return NotIntError
	}
	n += 1
	if err := c.db.Put(args[0], codec.EncodeStrVal(Int64ToSlice(n))); err != nil {
		return err
	}
	c.respWriter.writeInt(n)
	return nil
}

func cmdDecr(c *client, args Args) error {
	if !args.matchN(1) {
		return &WrongParamError{"decr"}
	}
	c.Lock(args.skey())
	defer c.Unlock(args.skey())

	value, err := c.db.Get(args[0])
	if err != nil {
		return err
	}
	if err := codec.CheckStrType(value); err != nil {
		return err
	}
	value = codec.DecodeStrKey(value)
	n, err := SliceToInt64(value)
	if err != nil {
		log.Println(err)
		return NotIntError
	}
	n -= 1
	if err := c.db.Put(args[0], codec.EncodeStrVal(Int64ToSlice(n))); err != nil {
		return err
	}
	c.respWriter.writeInt(n)
	return nil
}

func init() {
	register("get", cmdGet)
	register("set", cmdSet)
	register("incr", cmdIncr)
	register("decr", cmdDecr)
}
