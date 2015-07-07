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

package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/realglobe-Inc/go-lib/erro"
)

// テスト用の HTTP サーバー。
type HttpServer struct {
	*httptest.Server
	respCh chan *httpResponse
}

type httpResponse struct {
	status int
	header http.Header
	body   []byte

	reqCh chan<- *http.Request
}

func NewHttpServer(timeout time.Duration) (*HttpServer, error) {
	respCh := make(chan *httpResponse, 1024)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		failMsg := "Unknown error"
		defer func() {
			if failMsg != "" {
				http.Error(w, failMsg, http.StatusInternalServerError)
			}
		}()

		var resp *httpResponse
		select {
		case resp = <-respCh:
		case <-time.After(timeout):
			failMsg = "Timed out"
			return
		}

		// リクエストを渡す。
		if err := func() error {
			failed := true
			defer func() {
				if failed {
					resp.reqCh <- nil
				}
			}()

			// 抜けると Close されるので、置き換える。
			req := *r
			if buff, err := ioutil.ReadAll(r.Body); err != nil {
				return erro.Wrap(err)
			} else {
				req.Body = ioutil.NopCloser(bytes.NewReader(buff))
			}
			resp.reqCh <- &req

			failed = false
			return nil
		}(); err != nil {
			log.Warn(erro.Wrap(err))
			failMsg = erro.Unwrap(err).Error()
			return
		}

		failMsg = ""
		for key, vals := range resp.header {
			for _, val := range vals {
				w.Header().Add(key, val)
			}
		}
		w.WriteHeader(resp.status)
		if resp.body != nil {
			w.Write(resp.body)
		}
	})

	return &HttpServer{
		httptest.NewServer(mux),
		respCh,
	}, nil
}

// 次に返させるレスポンスを与える。
func (this *HttpServer) AddResponse(status int, header http.Header, body []byte) <-chan *http.Request {
	reqCh := make(chan *http.Request, 1)
	this.respCh <- &httpResponse{status, header, body, reqCh}
	return reqCh
}

// 待ち受けソケットを閉じる。
func (this *HttpServer) Close() {
	this.Server.Close()
	close(this.respCh)
}
