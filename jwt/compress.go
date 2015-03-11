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
	"compress/flate"
	"github.com/realglobe-Inc/go-lib/erro"
	"io/ioutil"
)

func defCompress(data []byte) ([]byte, error) {
	var buff bytes.Buffer
	if compressor, err := flate.NewWriter(&buff, flate.DefaultCompression); err != nil {
		return nil, erro.Wrap(err)
	} else if _, err := compressor.Write(data); err != nil {
		return nil, erro.Wrap(err)
	} else if compressor.Close(); err != nil {
		return nil, erro.Wrap(err)
	}
	return buff.Bytes(), nil
}

func defDecompress(data []byte) ([]byte, error) {
	decompressor := flate.NewReader(bytes.NewReader(data))
	defer decompressor.Close()
	buff, err := ioutil.ReadAll(decompressor)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return buff, nil
}
