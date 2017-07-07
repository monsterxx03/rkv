package db

import (
	rocksdb "github.com/tecbot/gorocksdb"
)

type DB struct {
	db *rocksdb.DB
}

func (d *DB) Close() {
	d.Close()
}

func (d *DB) Put(key, value []byte) error {
	wo := rocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	if err := d.db.Put(wo, key, value); err != nil {
		return err
	}
	return nil
}

func (d *DB) Get(key []byte) (*rocksdb.Slice, error) {
	ro := rocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	data, err := d.db.Get(ro, key)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *DB) Delete(key []byte) {}
func (d *DB) BatchWrite()       {}

func NewDB() *DB {
	bbto := rocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(rocksdb.NewLRUCache(3 << 30))
	opts := rocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := rocksdb.OpenDb(opts, "data")
	if err != nil {
		panic(err)
	}
	return &DB{db}
}

//func main() {
//	// open/create db
//	bbto := rocksdb.NewDefaultBlockBasedTableOptions()
//	bbto.SetBlockCache(rocksdb.NewLRUCache(3 << 30))
//	opts := rocksdb.NewDefaultOptions()
//	opts.SetBlockBasedTableFactory(bbto)
//	opts.SetCreateIfMissing(true)
//	db, err := rocksdb.OpenDb(opts, "data")
//	if err != nil {
//		panic(err)
//	}
//	// read/write options
//	ro := rocksdb.NewDefaultReadOptions()
//	wo := rocksdb.NewDefaultWriteOptions()
//	defer ro.Destroy()
//	defer wo.Destroy()
//	// fillin data
//	for i:=0; i< 10; i++ {
//		if err := db.Put(wo, []byte("foo" + strconv.Itoa(i)), []byte("bar")); err != nil {
//			panic(err)
//		}
//	}
//	// batch write
//	batch := rocksdb.NewWriteBatch()
//	batch.Put([]byte("b1"), []byte("v1"))
//	batch.Put([]byte("b2"), []byte("v2"))
//	batch.Delete([]byte("b1"))
//	db.Write(wo, batch)
//	// iter over
//	it := db.NewIterator(ro)
//	for it.SeekToFirst(); it.Valid(); it.Next() {
//		fmt.Println(string(it.Key().Data()), string(it.Value().Data()))
//	}
//
//	db.Close()
//}
