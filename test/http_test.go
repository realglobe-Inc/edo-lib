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
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	server, err := NewHttpServer(time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	reqCh := server.AddResponse(http.StatusOK, http.Header{"Test-Header": {"test header"}}, []byte("test body"))

	req, err := http.NewRequest("GET", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	rcvReq := <-reqCh

	if rcvReq.Method != req.Method {
		t.Fatal(rcvReq, req)
	} else if rcvReq.Host != req.Host {
		t.Fatal(rcvReq, req)
	} else if resp.StatusCode != http.StatusOK {
		t.Fatal(rcvReq, req)
	} else if resp.Header.Get("Test-Header") != "test header" {
		t.Error(rcvReq)
		t.Fatal(req)
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	} else if string(body) != "test body" {
		t.Error(string(body))
		t.Fatal("test body")
	}
}
