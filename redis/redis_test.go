package redis

import (
	"testing"
	"bufio"
	"bytes"
)

func Test_readLine(t *testing.T) {
	testData := []byte("1234567890")
	buf := bufio.NewReaderSize(bytes.NewReader(testData), 4096)
	data, err := readLine(buf)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, testData) {
		t.Error("Not equal")
	}
}

func Test_readLine_overbufer(t *testing.T) {
	testData := []byte("1234567890")
	buf := bufio.NewReaderSize(bytes.NewReader(testData), 5)
	data, err := readLine(buf)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, testData) {
		t.Error("Not equal")
	}
}