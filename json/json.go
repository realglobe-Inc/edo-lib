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

package json

import ()

func StringEscape(s string) string {
	return string(Escape([]byte(s)))
}

// http://www.json.org/index.html より。
func Escape(d []byte) []byte {
	output := []byte{}
	for _, r := range d {
		switch r {
		case '"':
			output = append(output, '\\', '"')
		case '\\':
			output = append(output, '\\', '\\')
		case '/':
			output = append(output, '\\', '/')
		case '\b':
			output = append(output, '\\', 'b')
		case '\f':
			output = append(output, '\\', 'f')
		case '\n':
			output = append(output, '\\', 'n')
		case '\r':
			output = append(output, '\\', 'r')
		case '\t':
			output = append(output, '\\', 't')
		default:
			output = append(output, r)
		}
	}
	return output
}
