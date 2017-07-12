/*
 string: key= key_name, value= TYPE_STR + len(value) + value + expire_at
 list:
   	- meta key: key=key_name, value = TYPE_LIST + size + expire-at
   	- one element one key: key= TYPE_LIST + key_name + seq, value=value
 set:
 		- meta key: key=key_name, value = TYPE_SET + size + expire_at
 		- one element one key: key = TYPE_SET + key_name + value, value = nil
 hash:
 		- meta key: key=key_name, value= TYPE_HASH + size + expire_at
 		- one field one key: key= TYPE_HASH + key_name +  field, value = value
 zset:
 		- meta key: key=key_name, value= TYPE_ZSET + size + expire_at
 		- one element one key: key= key_name + member, value = score
 		- one element one score key: TYPE_ZSET + key_name + score + member, value = nil

	expire_key: key_name + TYPE_EXPIRE, value: expire_at
*/

package codec

import (
	"encoding/binary"
	"errors"
)

const (
	StrType  byte = iota
	ListType
	HashType
	SetType
	ZSetType
)

func EncodeStrVal(value []byte) []byte {
	key := make([]byte, len(value)+1)
	key[0] = StrType
	copy(key[1:], value)
	return key
}

var WrongTypeError = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")

func DecodeStrKey(rawValue []byte) []byte {
	if len(rawValue) == 0 {
		return []byte{}
	}
	return rawValue[1:]
}

func DecodeType(rawValue []byte) byte {
	return rawValue[0]
}

func EncodeMetaVal(dataType byte, size int) []byte {
	buf := make([]byte, 1+4)
	buf[0] = dataType
	binary.BigEndian.PutUint32(buf[1:], uint32(size))
	return buf
}

func DecodeSize(rawValue []byte) uint32 {
	return binary.BigEndian.Uint32(rawValue[1:])
}

func EncodeListKey(keyName []byte, seq int) []byte {
	buf := make([]byte, 1+len(keyName)+4)
	buf[0] = ListType
	pos := 1
	copy(buf[pos:], keyName)
	pos += len(keyName)
	binary.BigEndian.PutUint32(buf[pos:], uint32(seq))
	return buf
}

func EncodeHashKey(keyName, fieldName []byte) []byte {
	buf := make([]byte, 1+len(keyName)+len(fieldName))
	buf[0] = HashType
	pos := 1
	pos += copy(buf[pos:], keyName)
	copy(buf[pos:], fieldName)
	return buf
}

func EncodeZSetKey(keyName, member []byte, score int64) []byte {
	buf := make([]byte, 1+len(keyName)+len(member)+8)
	buf[0] = ZSetType
	pos := 1
	pos += copy(buf[pos:], keyName)
	pos += copy(buf[pos:], member)
	binary.BigEndian.PutUint64(buf[pos:], uint64(score))
	return buf
}

func checkType(value []byte, dataType byte) error {
	if len(value) == 0 {
		return errors.New("Empty value")
	}
	if DecodeType(value) != dataType {
		return WrongTypeError
	}
	return nil
}

func CheckStrType(value []byte) error {
	return checkType(value, StrType)
}

func CheckListType(value []byte) error {
	return checkType(value, ListType)
}

func CheckHashType(value []byte) error {
	return checkType(value, HashType)
}

func CheckSetType(value []byte) error {
	return checkType(value, SetType)
}

func CheckZSetType(value []byte) error {
	return checkType(value, ZSetType)
}
