package server

import "strconv"

func Int64ToSlice(n int64) []byte {
	return strconv.AppendInt(nil, n, 10)
}

func SliceToInt64(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}


func FilterByte(data [][]byte, f func([]byte) bool) [][]byte {
	result := make([][]byte, 0)
	for _, v := range data {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
