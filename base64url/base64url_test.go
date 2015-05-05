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

package base64url

import (
	"bytes"
	"strings"
	"testing"
)

func TestBase64Url(t *testing.T) {
	for src := []byte{}; len(src) < 300; src = append(src, byte(len(src))) {
		if enc := EncodeToString(src); strings.Index(enc, "=") > 0 {
			t.Error(src)
			t.Fatal(enc)
		} else if dec, err := DecodeString(enc); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(dec, src) {
			t.Error(src)
			t.Error(enc)
			t.Fatal(dec)
		}
	}
}
