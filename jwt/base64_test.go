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

package jwt

import (
	"bytes"
	"strings"
	"testing"
)

func TestBase64Url(t *testing.T) {
	for b := []byte{}; len(b) < 100; b = append(b, byte(len(b))) {
		if s := base64UrlEncodeToString(b); strings.Index(s, "=") > 0 {
			t.Error(b)
			t.Error(s)
		} else if b2, err := base64UrlDecodeString(s); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(b2, b) {
			t.Error(b2)
			t.Error(b)
			t.Error(s)
		}
	}
}
