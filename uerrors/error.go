package uerrors

import (
	"errors"
	"strconv"
)

type Error struct {
	Err  error
	code int
}

func New(code int, message string) (err *Error) {
	err = new(Error)
	err.Err = errors.New(message)
	err.code = code
	return
}

func (err *Error) Error() (errStr string) {
	errStr = "code:" + strconv.Itoa(err.code) + " err:" + err.Err.Error()
	return
}
