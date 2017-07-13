package rocksdb

import (
	"github.com/monsterxx03/rkv/db/backend"
	"github.com/tecbot/gorocksdb"
	"github.com/monsterxx03/rkv/config"
)

type Backend struct {
}

func (s Backend) String() string {
	return "rocksdb"
}

func (s Backend) Open(cfg *config.Config) (backend.IDB, error) {
	db := newDB()
	if err := db.open(cfg); err != nil {
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

func (db *DB) open(cfg *config.Config) error {
	opts := db.initOptions(cfg)

	if rdb, err := gorocksdb.OpenDb(opts, cfg.RocksDB.DataDir); err != nil {
		return err
	} else {
		db.db = rdb
		return nil
	}
}

func (db *DB) initOptions(cfg *config.Config) *gorocksdb.Options {
	blo := gorocksdb.NewDefaultBlockBasedTableOptions()
	ropt := cfg.RocksDB

	blo.SetBlockCache(gorocksdb.NewLRUCache(ropt.BlockCache))
	blo.SetBlockSize(ropt.BlockSize)
	blo.SetFilterPolicy(gorocksdb.NewBloomFilter(ropt.BloomFilterBitsPerKey))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(blo)

	env := gorocksdb.NewDefaultEnv()
	env.SetBackgroundThreads(ropt.BackgroundThreads)
	env.SetHighPriorityBackgroundThreads(ropt.HighPriorityBackgroundThreads)
	opts.SetEnv(env)

	opts.SetCreateIfMissing(true)
	opts.SetCompression(gorocksdb.CompressionType(ropt.CompressionType))
	opts.SetWriteBufferSize(ropt.WriteBufferSize)
	opts.SetMaxWriteBufferNumber(ropt.MaxWriteBufferNumber)
	opts.SetMinWriteBufferNumberToMerge(ropt.MinWriteBufferNumberToMerge)
	opts.SetMaxOpenFiles(ropt.MaxOpenFiles)
	opts.SetNumLevels(ropt.NumLevels)
	opts.SetMaxBackgroundCompactions(ropt.MaxBackgroundCompactions)
	opts.SetMaxBackgroundFlushes(ropt.MaxBackgroundFlushes)
	opts.SetUseFsync(ropt.UseFsync)
	return opts
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

func (db *DB) MGet(keys [][]byte) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	values := make([][]byte, len(keys))
	for i, key := range keys {
		if val, err := db.Get(key); err != nil {
			return nil, err
		} else {
			values[i] = val
		}
	}
	return values, nil
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

func (db *DB) NewIter() backend.IIter {
	return db.db.NewIterator(db.defaultRo)
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
