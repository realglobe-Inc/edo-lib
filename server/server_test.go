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
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	logutil "github.com/realglobe-Inc/edo-lib/log"
	"github.com/realglobe-Inc/edo-lib/test"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

func init() {
	logutil.SetupConsole(logRoot, level.OFF)
}

type testParameter struct {
	port   int
	shutCh chan struct{}
}

func newTestParameter(port int) *testParameter {
	return &testParameter{
		port,
		make(chan struct{}, 10),
	}
}

func (this *testParameter) SocketPort() int {
	return this.port
}

func (this *testParameter) ShutdownChannel() chan struct{} {
	return this.shutCh
}

func TestServe(t *testing.T) {
	port, err := test.FreePort()
	if err != nil {
		t.Fatal(err)
	}
	param := newTestParameter(port)

	go func() {
		for {
			time.Sleep(time.Millisecond)

			r, _ := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+"/", nil)
			r.Header.Set("Connection", "close")
			resp, err := (&http.Client{}).Do(r)
			if err == nil {
				resp.Body.Close()
				break
			}
		}
		param.shutCh <- struct{}{}
	}()

	if err := Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), "tcp", "http", param); err != nil {
		t.Fatal(err)
	}
}

func TestServeRetry(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole(logRoot, level.ALL)
	// defer logutil.SetupConsole(logRoot, level.OFF)
	// ////////////////////////////////

	port, err := test.FreePort()
	if err != nil {
		t.Fatal(err)
	}
	param1 := newTestParameter(port)
	param2 := newTestParameter(port)

	// param1 で稼動させる。
	go Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), "tcp", "http", param1)

	// param1 で稼動しているのを確認する。
	for {
		time.Sleep(time.Millisecond)

		r, _ := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+"/", nil)
		r.Header.Set("Connection", "close")
		_, err := (&http.Client{}).Do(r)
		if err == nil {
			break
		}
	}

	body := []byte("abcde")
	go func() {
		// param2 で稼動しているのを確認できたら param2 を落とす。
		for {
			if func() bool {
				time.Sleep(time.Millisecond)

				r, _ := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+"/", nil)
				r.Header.Set("Connection", "close")
				resp, err := (&http.Client{}).Do(r)
				if err != nil {
					return false
				}
				defer resp.Body.Close()
				if buff, err := ioutil.ReadAll(resp.Body); err != nil {
					return false
				} else if bytes.Equal(buff, body) {
					return true
				}
				return false
			}() {
				break
			}
		}
		param2.shutCh <- struct{}{}
	}()

	go func() {
		// 少し経ったら param1 を落とす。
		time.Sleep(10 * time.Millisecond)
		param1.shutCh <- struct{}{}
	}()

	if err := Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(body) }), "tcp", "http", param2); err != nil {
		t.Fatal(err)
	}
}
