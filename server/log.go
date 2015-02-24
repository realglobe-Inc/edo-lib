package server

import (
	"github.com/realglobe-Inc/go-lib/rglog"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net/http"
	"net/http/httputil"
)

var log = rglog.Logger("github.com/realglobe-Inc/edo/util/http")

func LogRequest(lv level.Level, r *http.Request, useBody bool, args ...interface{}) {
	if log.IsLoggable(lv) {
		buff, _ := httputil.DumpRequest(r, useBody)
		a := []interface{}{"Request: "}
		if len(args) > 0 {
			a = append(a, args...)
			a = append(a, ": ")
		}
		a = append(a, string(buff))
		log.Log(lv, a...)
	}
}

func LogResponse(lv level.Level, r *http.Response, useBody bool, args ...interface{}) {
	if log.IsLoggable(lv) {
		buff, _ := httputil.DumpResponse(r, useBody)
		a := []interface{}{"Response: "}
		if len(args) > 0 {
			a = append(a, args...)
			a = append(a, ": ")
		}
		a = append(a, string(buff))
		log.Log(lv, a...)
	}
}
