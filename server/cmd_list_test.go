package server

import (
	"testing"
	"github.com/monsterxx03/rkv/codec"
)

func TestCmdLpush(t *testing.T) {
	conn := GetTestConn()
	key := "testL"
	conn.Set("testStr", "abc",0)
	// lpush on non-list key
	if _, err := conn.LPush("testStr", 1).Result(); err.Error() != codec.WrongTypeError.Error() {
		t.Fatal(err)
	}
	// lpush on non-exist key
	if value, err := conn.LPush(key, 1).Result(); err != nil {
		t.Fatal(err)
	} else if value != 1 {
		t.Fatal(value)
	}
	// lpush on existing key
	if value, err := conn.LPush(key, "a", "b").Result(); err != nil {
		t.Fatal(err)
	} else if value != 3 {
		t.Fatal(value)
	}
}


func TestCmdLpop(t *testing.T) {
}