package server

import (
	"bufio"
	logutil "github.com/realglobe-Inc/edo-toolkit/util/log"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net/http"
	"strings"
	"testing"
)

func TestLogRequest(t *testing.T) {
	logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)

	req, err := http.NewRequest("GET", "http://example.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	LogRequest(level.ALL, req, true, "test")
}

func TestLogResponse(t *testing.T) {
	logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)

	resp, err := http.ReadResponse(bufio.NewReader(strings.NewReader(`HTTP/1.1 200 OK
Server: nginx/1.7.9
Date: Tue, 10 Feb 2015 02:46:19 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 1
Connection: keep-alive

a
`)), nil)
	if err != nil {
		t.Fatal(err)
	}

	LogResponse(level.ALL, resp, true, "test")
}
