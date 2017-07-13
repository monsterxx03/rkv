package server

import (
	"github.com/monsterxx03/rkv/codec"
	"github.com/monsterxx03/rkv/db/backend"
)

func delList(batch backend.IBatch, key, metaValue []byte) error {
	// size := int(codec.DecodeSize(metaValue))
	return nil
}

func cmdLpush(c *client, args Args) error {
	if len(args) < 2 {
		return &WrongParamError{"lpush"}
	}
	c.Lock(args.skey())
	defer c.Unlock(args.skey())

	// value is metakey's value
	value, err := c.db.Get(args.key())
	if err != nil {
		return err
	}
	size := 0
	if len(value) > 0 {
		if err := codec.CheckListType(value); err != nil {
			return err
		}
		size = int(codec.DecodeSize(value))
	}
	value = codec.EncodeMetaVal(codec.ListType, int(size) + len(args.values()))
	// write meta key
	batch := c.db.NewBatch()
	batch.Put(args.key(), value)
	for i, v:= range args[1:] {
		batch.Put(codec.EncodeListKey(args.key(), size + i), v)
	}
	batch.Commit()
	c.respWriter.writeInt(int64(size + len(args.values())))
	return nil
}

func cmdLpop(c *client, args Args) error {
	if len(args) != 1 {
		return &WrongParamError{"lpop"}
	}
	return nil
}

func cmdLrange(c *client, args Args) error {
	if len(args) != 3 {
		return &WrongParamError{"lrange"}
	}
	return nil
}

func cmdRpush(c *client, args Args) error {
	if len(args) == 0 {
		return &WrongParamError{"rpush"}
	}
	c.Lock(args.skey())
	defer c.Unlock(args.skey())
	return nil
}

func cmdRpop(c *client, args Args) error {
	if len(args) != 1 {
		return &WrongParamError{"rpop"}
	}
	return nil
}



func init() {
	register("lpush", cmdLpush)
	register("lpop", cmdLpop)
	register("lrange", cmdLrange)
	register("rpush", cmdRpush)
	register("rpop", cmdRpop)
}