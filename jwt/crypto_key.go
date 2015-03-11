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
	"crypto/aes"
	"github.com/realglobe-Inc/go-lib/erro"
)

var aesKwInitVec = []byte{0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6, 0xa6}

// AES Key Wrap 暗号。
func encryptAesKw(key []byte, plain []byte) (encrypted []byte, err error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	n := len(plain) / 8
	ari := make([]byte, 16)
	copy(ari[:8], aesKwInitVec)
	r := make([]byte, len(plain))
	copy(r, plain)

	b := make([]byte, 16)
	for j := 0; j <= 5; j++ {
		for i := 1; i <= n; i++ {
			ri := r[8*(i-1) : 8*i]
			copy(ari[8:], ri)
			cip.Encrypt(b, ari)
			_xorBytes(ari[:8], b[:8], int64(n*j+i))
			copy(ri, b[8:])
		}
	}

	c := make([]byte, len(r)+8)
	copy(c[:8], ari[:8])
	copy(c[8:], r)

	return c, nil
}

// AES Key Wrap 復号。
func decryptAesKw(key []byte, encrypted []byte) (plain []byte, err error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	n := len(encrypted)/8 - 1
	ari := make([]byte, 16)
	copy(ari[:8], encrypted[:8])
	r := make([]byte, len(encrypted)-8)
	copy(r, encrypted[8:])

	b := make([]byte, 16)
	for j := 5; j >= 0; j-- {
		for i := n; i >= 1; i-- {
			ri := r[8*(i-1) : 8*i]
			_xorBytes(ari[:8], ari[:8], int64(n*j+i))
			copy(ari[8:], ri)
			cip.Decrypt(b, ari)
			copy(ari[:8], b[:8])
			copy(ri, b[8:])
		}
	}

	if !bytes.Equal(ari[:8], aesKwInitVec) {
		return nil, erro.New("verification error")
	}

	return r, nil
}

func _xorBytes(dst, src1 []byte, v int64) {
	for i := 7; i >= 0; i-- {
		dst[i] = src1[i] ^ byte(v)
		v >>= 8
	}
}
