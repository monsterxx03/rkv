package rocksdb

import (
	"github.com/monsterxx03/rkv/db/backend"
	"github.com/tecbot/gorocksdb"
)

type Backend struct {
}

func (s Backend) String() string {
	return "rocksdb"
}

func (s Backend) Open() (backend.IDB, error) {
	db := newDB()
	if err := db.open(); err != nil {
		return nil, err
	}
	return db, nil
}

type DB struct {
	db *gorocksdb.DB
	defaultWo *gorocksdb.WriteOptions
	defaultRo *gorocksdb.ReadOptions
}

// WriteBatch wrapper
type WriteBatch struct {
	db *DB
	wb *gorocksdb.WriteBatch
}

func (b *WriteBatch) Close() {
	b.wb.Destroy()
}

func (b *WriteBatch) Put(key, value []byte) {
	b.wb.Put(key, value)
}

func (b *WriteBatch) Delete(key []byte) {
	b.wb.Delete(key)
}

func (b *WriteBatch) Commit() error {
	if err := b.db.db.Write(b.db.defaultWo, b.wb); err != nil {
		return err
	}
	return nil
}

func newDB() *DB {
	rdb := new(DB)
	rdb.defaultWo = gorocksdb.NewDefaultWriteOptions()
	rdb.defaultRo = gorocksdb.NewDefaultReadOptions()
	return rdb
}

func (db *DB) open() error {
	// TODO init options dynamically from config file
	blo := gorocksdb.NewDefaultBlockBasedTableOptions()
	blo.SetBlockCache(gorocksdb.NewLRUCache(1073741824))
	blo.SetBlockSize(65536)
	blo.SetFilterPolicy(gorocksdb.NewBloomFilter(10))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(blo)

	env := gorocksdb.NewDefaultEnv()
	env.SetBackgroundThreads(16)
	env.SetHighPriorityBackgroundThreads(1)
	opts.SetEnv(env)

	opts.SetCreateIfMissing(true)
	opts.SetCompression(gorocksdb.NoCompression)
	opts.SetWriteBufferSize(134217728)
	opts.SetMaxWriteBufferNumber(6)
	opts.SetMinWriteBufferNumberToMerge(2)
	opts.SetMaxOpenFiles(1024)
	opts.SetNumLevels(7)
	opts.SetMaxBackgroundCompactions(15)
	opts.SetMaxBackgroundFlushes(1)
	opts.SetUseFsync(false)

	if rdb, err := gorocksdb.OpenDb(opts, "data"); err != nil {
		return err
	} else {
		db.db = rdb
		return nil
	}
}

func (db *DB) Close() error {
	db.db.Close()
	return nil
}

func (db *DB) Put(key, value []byte) error {
	if err := db.db.Put(db.defaultWo, key, value); err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	if result, err := db.db.Get(db.defaultRo, key); err != nil {
		return nil, err
	} else {
		return result.Data(), nil
	}
}

func (db *DB) Delete(key []byte) error {
	if err := db.db.Delete(db.defaultWo, key); err != nil {
		return err
	}
	return nil
}

func (db *DB) NewBatch() backend.IBatch {
	wb := new(WriteBatch)
	wb.wb = gorocksdb.NewWriteBatch()
	wb.db = db
	return wb
}

func init() {
	backend.RegisterBackend(Backend{})
}

/*
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

*/
