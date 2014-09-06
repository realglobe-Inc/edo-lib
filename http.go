package util

import (
	"net/http"
	"net/http/httputil"
)

const ContentTypeJson string = "application/json"

func LogRequest(r *http.Request, useBody bool) {
	buff, _ := httputil.DumpRequest(r, useBody)
	log.Debug("Request: " + string(buff))
}

func LogResponse(r *http.Response, useBody bool) {
	buff, _ := httputil.DumpResponse(r, useBody)
	log.Debug("Response: " + string(buff))
}
