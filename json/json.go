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
	output := ""
	for _, r := range s {
		switch r {
		case '"':
			output += "\\\""
		case '\\':
			output += "\\\\"
		case '/':
			output += "\\/"
		case '\n':
			output += "\\n"
		case '\r':
			output += "\\r"
		case '\t':
			output += "\\t"
		case '\b':
			output += "\\b"
		case '\f':
			output += "\\f"
		default:
			output += string(r)
		}
	}
	return output
}
