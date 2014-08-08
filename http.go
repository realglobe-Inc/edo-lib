package util

import (
	"net/http"
	"net/http/httputil"
)

const ContentTypeJson string = "application/json"

func LogRequest(r *http.Request, body bool) {
	buff, _ := httputil.DumpRequest(r, body)
	log.Debug("Request: " + string(buff))
}

func LogResponse(r *http.Response, body bool) {
	buff, _ := httputil.DumpResponse(r, true)
	log.Debug("Response: " + string(buff))
}
