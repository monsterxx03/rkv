package server

import "strconv"

func Int64ToSlice(n int64) []byte {
	return strconv.AppendInt(nil, n, 10)
}
