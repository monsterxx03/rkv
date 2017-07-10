package server

import (
	"fmt"
)

type WrongParamError struct {
	cmd string
}

func (e *WrongParamError) Error() string {
	return fmt.Sprintf("wrong number of arguments for '%s' command",  e.cmd)
}
