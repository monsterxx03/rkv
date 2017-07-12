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

var NotIntError = errors.New("ERR value is not an integer or out of range")
