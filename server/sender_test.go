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
	"net/http"
	"testing"
)

func TestParseSender(t *testing.T) {
	r, err := http.NewRequest("GET", "https://server.example.org/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.RemoteAddr = "192.168.0.18:55555"
	if src := ParseSender(r); src != "192.168.0.18:55555" {
		t.Error(src)
		t.Fatal("192.168.0.18:55555")
	}

	r.Header.Set("X-Forwarded-For", "203.0.113.34, 192.168.0.12")
	if src := ParseSender(r); src != "203.0.113.34" {
		t.Error(src)
		t.Fatal("203.0.113.34")
	}
}
