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

package password

import (
	"bytes"
	"testing"
)

func TestCalculate(t *testing.T) {
	if hash, err := Calculate("pbkdf2:sha256:1000", test_salt, test_passwd); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(hash, test_pbkdf2Hash) {
		t.Error(hash)
		t.Fatal(test_pbkdf2Hash)
	}
}

func TestCalculateError(t *testing.T) {
	if _, err := Calculate("", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("unknown", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	}
}

func TestPbkdf2CalculateError(t *testing.T) {
	if _, err := Calculate("pbkdf2:unknown:1000", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256:invalid", test_salt, test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256:1000"); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256:1000", test_salt); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256:1000", "abcde", test_passwd); err == nil {
		t.Fatal("no error")
	} else if _, err := Calculate("pbkdf2:sha256:1000", test_salt, []byte{0, 1, 2, 3, 4}); err == nil {
		t.Fatal("no error")
	}
}
