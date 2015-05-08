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
	"testing"
)

func TestRedisServer(t *testing.T) {
	server, err := NewRedisServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	if server.Address() == "" {
		t.Fatal("no address")
	} else if server.Pool() == nil {
		t.Fatal("no pool")
	}
}
