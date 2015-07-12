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

package reader

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/realglobe-Inc/edo-lib/prand"
)

// メモリだけでテスト。
func TestResettableMemory(t *testing.T) {
	data := []byte(prand.New(time.Hour).Bytes(100))
	base := ioutil.NopCloser(bytes.NewReader(data))
	buff := NewResettable(base, len(data), test_pref, 0)
	defer buff.Dispose()

	for i := 0; i <= len(data); i++ {
		buff.Reset()
		data2 := make([]byte, i)
		if _, err := io.ReadFull(buff, data2); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(data2, data[:len(data2)]) {
			t.Error(i)
			t.Error(string(data2))
			t.Fatal(string(data[:len(data2)]))
		}
	}
	if _, err := buff.Read(make([]byte, 1)); err != io.EOF {
		t.Fatal(err)
	}
}

// ファイルもテスト。
func TestResettable(t *testing.T) {
	data := []byte(prand.New(time.Hour).String(100))
	base := ioutil.NopCloser(bytes.NewReader(data))
	buff := NewResettable(base, len(data)/2, test_pref, (len(data)+1)/2)
	defer buff.Dispose()

	for i := 0; i <= len(data); i++ {
		buff.Reset()
		data2 := make([]byte, i)
		if _, err := io.ReadFull(buff, data2); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(data2, data[:len(data2)]) {
			t.Error(i)
			t.Error(string(data2))
			t.Fatal(string(data[:len(data2)]))
		}
	}
	if _, err := buff.Read(make([]byte, 1)); err != io.EOF {
		t.Fatal(err)
	}
}

func TestResettableLast(t *testing.T) {
	data := []byte(prand.New(time.Hour).String(100))
	base := ioutil.NopCloser(bytes.NewReader(data))
	buff := NewResettable(base, len(data)/2, test_pref, (len(data)+1)/2)
	defer buff.Dispose()

	data2 := make([]byte, len(data)/2)
	if _, err := io.ReadFull(buff, data2); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(data2, data[:len(data2)]) {
		t.Error(string(data2))
		t.Fatal(string(data[:len(data2)]))
	}
	buff.LastReset()

	data3 := make([]byte, len(data))
	if _, err := io.ReadFull(buff, data3); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(data3, data) {
		t.Error(string(data3))
		t.Fatal(string(data))
	}

	if _, err := buff.Read(make([]byte, 1)); err != io.EOF {
		t.Fatal(err)
	}

	if buff.file != nil {
		t.Fatal("file opened")
	}
}

// いっぱい読んだらリセットできなくなるけど貯めないことのテスト。
func TestOverflow(t *testing.T) {
	data := []byte(prand.New(time.Hour).String(100))
	base := ioutil.NopCloser(bytes.NewReader(data))
	buff := NewResettable(base, len(data)/4, test_pref, len(data)/4)
	defer buff.Dispose()

	data2 := make([]byte, len(data))
	for i := 0; i < 10; i++ {
		if _, err := io.ReadFull(buff, data2[i*10:i*10+len(data2)/10]); err != nil {
			t.Fatal(err)
		}
	}
	if !bytes.Equal(data2, data) {
		t.Error(string(data2))
		t.Fatal(string(data))
	}

	if _, err := buff.Read(make([]byte, 1)); err != io.EOF {
		t.Fatal(err)
	}

	buff.fileW.Flush()
	if fi, _ := buff.file.Stat(); fi.Size() == int64(len(data)) {
		t.Fatal("all data reserved")
	}
}
