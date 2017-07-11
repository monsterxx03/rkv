package server

import (
	"bufio"
	"errors"
	"strconv"
	"io"
	"fmt"
	"bytes"
	"net"
)

//RESP: Redis Serialization Protocol
const (
	respSimpleString = '+'
	respERROR        = '-'
	respInt          = ':'
	respString       = '$'
	respArray        = '*'
)
var DELIMS  = []byte("\r\n")

type RESPReader struct {
	buf *bufio.Reader
}

type RESPWriter struct {
	buf *bufio.Writer
}

func (reader *RESPReader) ParseRequest() ([][]byte, error) {
	headerByte, err := reader.buf.ReadByte()
	if err != nil {
		return nil, err
	}
	reader.buf.UnreadByte()
	if headerByte != respArray { // inline command
		return parseInlineCmd(reader.buf)
	}
	line, err := readLine(reader.buf)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("Empty RESP data")
	}
	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}
	cmds := make([][]byte, length)
	for i := range cmds {
		cmds[i], err = parseBulkString(reader.buf)
		if err != nil {
			return nil, err
		}
	}
	return cmds, nil
}

func parseInlineCmd(buf *bufio.Reader) ([][]byte, error) {
	line, err := buf.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	r := make([][]byte, 1)
	scanner := bufio.NewScanner(bytes.NewReader(line))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		b := scanner.Bytes()
		fmt.Println(string(b))
		r = append(r, b)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return r, nil
}

func parseBulkString(buf *bufio.Reader) ([]byte, error) {
	line, err := readLine(buf)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("Empty bulk string")
	}
	if line[0] != respString {
		return nil, errors.New("Invalid bulk string")
	}
	length, err := strconv.Atoi(string(line[1:]))
	if length < 0 || err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(buf, data); err != nil {
		return nil, err
	}
	// consume last \r\n
	if line, err := readLine(buf); err != nil {
		return nil, err
	} else if len(line) != 0 {
		return nil, errors.New("Invalid bulk string")
	}
	return data, nil
}

func NewRESPReader(reader *bufio.Reader) *RESPReader {
	return &RESPReader{reader}
}

func readLine(buf *bufio.Reader) ([]byte, error) {
	data, isPrefix, err := buf.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		_data, err := readLine(buf)
		if err != nil {
			return nil, err
		}
		return append(data, _data...), nil
	}
	return data, nil
}


func NewRESPWriter(conn net.Conn, size int) *RESPWriter {
	return &RESPWriter{buf: bufio.NewWriterSize(conn, size)}
}

func (w *RESPWriter) flush() {
	w.buf.Flush()
}

func (w *RESPWriter) writeError(err error) {
	w.buf.Write([]byte("-"))
	if err != nil {
		w.buf.Write([]byte(err.Error()))
	}
	w.buf.Write(DELIMS)
}

// write simple string to response
func (w *RESPWriter) writeStr(s string) {
	w.buf.Write([]byte("+"))
	w.buf.Write([]byte(s))
	w.buf.Write(DELIMS)
}
