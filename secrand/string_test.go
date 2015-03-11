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

package secrand

import (
	"testing"
)

func TestString(t *testing.T) {
	for i := 0; i < 100; i++ {
		buff, err := String(i)
		if err != nil {
			t.Fatal(err)
		} else if len(buff) != i {
			t.Error(i, len(buff), " "+buff)
		} else if len(buff) > 0 && buff[len(buff)-1] == '=' {
			t.Error(i, len(buff), " "+buff)
		}
	}
}
