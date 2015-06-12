// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bufio"
	logutil "github.com/realglobe-Inc/edo-lib/log"
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

	LogRequest(level.ALL, req, true, "test: ")
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

	LogResponse(level.ALL, resp, true, "test: ")
}
