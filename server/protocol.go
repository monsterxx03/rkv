package server

import (
	"bufio"
	"errors"
	"strconv"
	"io"
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

var (
	DELIMS = []byte("\r\n")
)

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
	length, err := parseLen(line[1:])
	if length < 0 || err != nil {
		return nil, err
	}
	data := make([][]byte, length)
	for i := range data {
		data[i], err = parseBulkString(reader.buf)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func parseInlineCmd(buf *bufio.Reader) ([][]byte, error) {
	line, err := buf.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	r := make([][]byte, 1)
	scanner := bufio.NewScanner(bytes.NewReader(line))
	scanner.Split(bufio.ScanWords)
	ops := 0
	for scanner.Scan() {
		b := scanner.Bytes()
		if ops >= 1 {
			r = append(r, b)
		} else {
			r[ops] = b
			ops++
		}
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
	length, err := parseLen(line[1:])
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

func parseLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errors.New("malformed length")
	}

	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, errors.New("illegal bytes in length")
		}
		n += int(b - '0')
	}

	return n, nil
}

func NewRESPWriter(conn net.Conn, size int) *RESPWriter {
	return &RESPWriter{buf: bufio.NewWriterSize(conn, size)}
}

func (w *RESPWriter) flush() error {
	return w.buf.Flush()
}

// return error message to client
func (w *RESPWriter) writeError(err error) {
	w.buf.WriteRune(respERROR)
	if err != nil {
		w.buf.WriteString(err.Error())
	}
	w.buf.Write(DELIMS)
}

// return simple string to client
func (w *RESPWriter) writeStr(s string) {
	w.buf.WriteRune(respSimpleString)
	w.buf.WriteString(s)
	w.buf.Write(DELIMS)
}

// return bulk string to client
func (w *RESPWriter) writeBulkStr(s []byte) {
	w.buf.WriteRune(respString)
	if len(s) > 0 {
		w.buf.WriteString(strconv.Itoa(len(s)))
		w.buf.Write(DELIMS)
		w.buf.Write(s)
		w.buf.Write(DELIMS)
	} else {
		w.buf.WriteString("-1")
		w.buf.Write(DELIMS)
		return
	}
}

func (w *RESPWriter) writeInt(n int64) {
	w.buf.WriteRune(respInt)
	w.buf.Write(Int64ToSlice(n))
	w.buf.Write(DELIMS)
}
