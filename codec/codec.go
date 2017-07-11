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

func DecodeStrKey(rawValue []byte) []byte {
	return rawValue[1:]
}

func DecodeType(rawValue []byte) byte {
	return rawValue[0]
}

func EncodeMetaKey(keyName []byte, dataType byte, size int32) []byte {
	buf := make([]byte, len(keyName)+1)
	buf[0] = dataType
	copy(buf, keyName)
	binary.BigEndian.PutUint32(buf, uint32(size))
	return buf
}

func EncodeListKey(keyName []byte, seq int32) []byte {
	buf := make([]byte, len(keyName)+1)
	buf[0] = ListType
	copy(buf, keyName)
	binary.BigEndian.PutUint32(buf, uint32(seq))
	return buf
}

func EncodeHashKey(keyName, fieldName []byte) []byte {
	buf := make([]byte, len(keyName)+1+len(fieldName))
	buf[0] = HashType
	pos := 1
	pos += copy(buf, keyName)
	copy(buf, fieldName)
	return buf
}
