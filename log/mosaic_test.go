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

package log

import (
	"testing"
)

func TestMosaic(t *testing.T) {
	for s := ""; len(s) < 2*DisplayLength; s += string('a' + len(s)%26) {
		if len(Mosaic(s)) > DisplayLength {
			t.Error(s)
			t.Fatal(Mosaic(s))
		}
	}
}
