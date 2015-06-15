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

// 末尾に = を足さない Base64URL エンコード。
package base64url

import (
	"encoding/base64"
	"github.com/realglobe-Inc/go-lib/erro"
)

func Decode(src []byte) ([]byte, error) {
	rem := len(src) % 4

	switch rem {
	case 2:
		// 上書きしてしまう場合は戻す。
		if cap(src) > len(src)+2 {
			buff := src[:len(src)+2]
			old := []byte{buff[len(src)], buff[len(src)+1]}
			defer func() {
				src[len(src)-2] = old[0]
				src[len(src)-1] = old[1]
			}()
		}
		src = append(src, '=', '=')
	case 3:
		// 上書きしてしまう場合は戻す。
		if cap(src) > len(src)+1 {
			buff := src[:len(src)+1]
			old := []byte{buff[len(src)]}
			defer func() {
				src[len(src)-1] = old[0]
			}()
		}
		src = append(src, '=')
	}

	dst := make([]byte, base64.URLEncoding.DecodedLen(len(src)))
	n, err := base64.URLEncoding.Decode(dst, src)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return dst[:n], nil
}

func DecodeString(s string) ([]byte, error) {
	return Decode([]byte(s))
}

func Encode(src []byte) []byte {
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(dst, src)

	switch len(src) % 3 {
	case 1:
		dst = dst[:len(dst)-2]
	case 2:
		dst = dst[:len(dst)-1]
	}

	return dst
}

func EncodeToString(src []byte) string {
	return string(Encode(src))
}
