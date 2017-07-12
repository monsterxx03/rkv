package server

import "strconv"

func Int64ToSlice(n int64) []byte {
	return strconv.AppendInt(nil, n, 10)
}

func SliceToInt64(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}
