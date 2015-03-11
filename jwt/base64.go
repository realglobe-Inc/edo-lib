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
	"encoding/base64"
)

// 末尾に = を足さない Base64URL エンコード。

func base64UrlDecodeString(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func base64UrlEncode(src []byte) []byte {
	buff := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(buff, src)
	switch len(src) % 3 {
	case 1:
		return buff[:len(buff)-2]
	case 2:
		return buff[:len(buff)-1]
	default:
		return buff
	}
}

func base64UrlEncodeToString(src []byte) string {
	return string(base64UrlEncode(src))
}

func Base64UrlDecodeString(s string) ([]byte, error) { return base64UrlDecodeString(s) }
func Base64UrlEncode(src []byte) []byte              { return base64UrlEncode(src) }
func Base64UrlEncodeToString(src []byte) string      { return base64UrlEncodeToString(src) }
