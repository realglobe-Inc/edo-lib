package server

import (
	"fmt"
)

// HTTP のステータスコードを付加したエラー。
type StatusError struct {
	status int
	msg    string

	cause error
}

func NewStatusError(status int, msg string, cause error) *StatusError {
	return &StatusError{status, msg, cause}
}

func (err *StatusError) Error() string {
	buff := err.msg
	if err.cause != nil {
		buff += fmt.Sprintln()
		buff += "caused by: "
		buff += err.cause.Error()
	}
	return buff
}

func (err *StatusError) Status() int {
	return err.status
}

func (err *StatusError) Message() string {
	return err.msg
}

func (err *StatusError) Cause() error {
	return err.cause
}
