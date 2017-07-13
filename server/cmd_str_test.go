package server

import (
	"testing"
	"github.com/go-redis/redis"
)

func TestCmdSet(t *testing.T) {
	conn := GetTestConn()
	key := "test1"
	value := "abc"
	intKey :=  "int1"
	// should return nil
	if val, err := conn.Get(key).Result(); err != redis.Nil {
		t.Fatal(err)
	} else if val != "" {
		t.Fatal("not empty resp: ", val)
	}

	// test set
	if val, err := conn.Set(key, value, 0).Result(); err != nil {
		t.Fatal(err)
	} else if val != "OK" {
		t.Fatal("bad resp: ", val)
	}

	// test get
	if val, err := conn.Get(key).Result(); err != nil {
		t.Fatal(err)
	} else if val != value {
		t.Fatal("bad resp: ", val)
	}

	// test incr on string value
	if val, err := conn.Incr(key).Result(); err.Error() != NotIntError.Error() {
		t.Fatal(err, val)
	}
	// test incr
	if val, err := conn.Incr(intKey).Result(); err != nil {
		t.Fatal(err)
	} else if val != 1 {
		t.Fatal("bad resp: ", val)
	}
	conn.Incr(intKey)
	if val, err := conn.Get(intKey).Result(); err != nil {
		t.Fatal(err)
	} else if val != "2" {
		t.Fatal("bad resp: ", val)
	}

	// test decr
	conn.Decr(intKey)
	conn.Decr(intKey)
	if val, err := conn.Decr(intKey).Result(); err != nil {
		t.Fatal(err)
	} else if val != -1 {
		t.Fatal("Bad resp: ", val)
	}
	if val, err := conn.Get(intKey).Result(); err != nil {
		t.Fatal(err)
	} else if val != "-1" {
		t.Fatal("Bad resp", val)
	}
}
