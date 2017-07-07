package store

import (
	"strconv"
	rocksdb "github.com/tecbot/gorocksdb"
	"fmt"
)

func main() {
	// open/create db
	bbto := rocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(rocksdb.NewLRUCache(3 << 30))
	opts := rocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := rocksdb.OpenDb(opts, "data")
	if err != nil {
		panic(err)
	}
	// read/write options
	ro := rocksdb.NewDefaultReadOptions()
	wo := rocksdb.NewDefaultWriteOptions()
	defer ro.Destroy()
	defer wo.Destroy()
	// fillin data
	for i:=0; i< 10; i++ {
		if err := db.Put(wo, []byte("foo" + strconv.Itoa(i)), []byte("bar")); err != nil {
			panic(err)
		}
	}
	// batch write
	batch := rocksdb.NewWriteBatch()
	batch.Put([]byte("b1"), []byte("v1"))
	batch.Put([]byte("b2"), []byte("v2"))
	batch.Delete([]byte("b1"))
	db.Write(wo, batch)
	// iter over
	it := db.NewIterator(ro)
	for it.SeekToFirst(); it.Valid(); it.Next() {
		fmt.Println(string(it.Key().Data()), string(it.Value().Data()))
	}

	db.Close()
}
