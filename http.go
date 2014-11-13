package util

import (
	"fmt"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net/http"
	"net/http/httputil"
)

const (
	ContentTypeJson string = "application/json"
	ContentTypeHtml        = "text/html"
)

func LogRequest(r *http.Request, useBody bool) {
	if log.IsLoggable(level.DEBUG) {
		buff, _ := httputil.DumpRequest(r, useBody)
		log.Debug("Request: " + string(buff))
	}
}

func LogResponse(r *http.Response, useBody bool) {
	if log.IsLoggable(level.DEBUG) {
		buff, _ := httputil.DumpResponse(r, useBody)
		log.Debug("Response: " + string(buff))
	}
}

// HTTP のステータスコードを付加したエラー。
type HttpStatusError struct {
	status int
	msg    string

	cause error
}

func NewHttpStatusError(status int, msg string, cause error) error {
	return &HttpStatusError{status, msg, cause}
}

func (err *HttpStatusError) Error() string {
	buff := err.msg
	if err.cause != nil {
		buff += fmt.Sprintln()
		buff += "caused by: "
		buff += err.cause.Error()
	}
	return buff
}

func (err *HttpStatusError) Status() int {
	return err.status
}

func (err *HttpStatusError) Message() string {
	return err.msg
}

func (err *HttpStatusError) Cause() error {
	return err.cause
}
