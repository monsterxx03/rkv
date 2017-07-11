package server

import (
	"fmt"
	"errors"
)

type WrongParamError struct {
	cmd string
}

func (e *WrongParamError) Error() string {
	return fmt.Sprintf("Err wrong number of arguments for '%s' command",  e.cmd)
}

var WrongTypeError = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
