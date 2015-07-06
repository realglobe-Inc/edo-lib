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

// 暗号関係。
package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/realglobe-Inc/go-lib/erro"
)

func perseBlock(block *pem.Block) (interface{}, error) {
	var key interface{}
	var err error
	switch block.Type {
	case "PUBLIC KEY":
		key, err = x509.ParsePKIXPublicKey(block.Bytes)
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return key, nil
}

// PEM 形式のデータから最初の鍵を取り出す。
func ParsePem(data []byte) (interface{}, error) {
	for block, rest := pem.Decode(data); block != nil; block, rest = pem.Decode(rest) {
		if key, err := perseBlock(block); err != nil {
			return nil, erro.Wrap(err)
		} else if key != nil {
			return key, nil
		}
	}
	return nil, erro.New("no supported key")
}

// PEM 形式のファイルから最初の鍵を読む。
func ReadPem(path string) (interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return ParsePem(data)
}

// PEM 形式のデータから鍵を取り出す。
func ParsePemAll(data []byte) ([]interface{}, error) {
	keys := []interface{}{}
	for block, rest := pem.Decode(data); block != nil; block, rest = pem.Decode(rest) {
		if key, err := perseBlock(block); err != nil {
			return nil, erro.Wrap(err)
		} else if key != nil {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// PEM 形式のファイルから鍵を読む。
func ReadPemAll(path string) ([]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return ParsePemAll(data)
}
