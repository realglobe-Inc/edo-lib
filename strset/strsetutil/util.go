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

// 文字列集合関係。
package strsetutil

// 文字列集合をつくる。
func New(strs ...string) map[string]bool {
	s := map[string]bool{}
	for _, str := range strs {
		s[str] = true
	}
	return s
}

// s1 が s2 を含むかどうか。
func Contains(s1, s2 map[string]bool) bool {
	for k := range s2 {
		if !s1[k] {
			return false
		}
	}
	return true
}
