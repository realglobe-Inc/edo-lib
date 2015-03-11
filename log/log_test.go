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

package log

import (
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net"
	"os"
	"path/filepath"
	"testing"
)

const testLabel = "edo-test"

// テストしたかったら fluentd サーバーを立ててから。
var testFluAddr = "localhost:24224"

func init() {
	if testFluAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := net.Dial("tcp", testFluAddr)
		if err != nil {
			testFluAddr = ""
		} else {
			conn.Close()
		}
	}
}

type testOption struct{}

func (opt *testOption) LogFilePath() string       { return filepath.Join(os.TempDir(), testLabel+".log") }
func (opt *testOption) LogFileLimit() int64       { return 10 * (1 << 20) }
func (opt *testOption) LogFileNumber() int        { return 10 }
func (opt *testOption) LogFluentdAddress() string { return testFluAddr }
func (opt *testOption) LogFluentdTag() string     { return testLabel }

var testOpt = &testOption{}

func TestSetup(t *testing.T) {
	for _, lv := range level.Values() {
		for _, typ := range []string{TypeConsole, TypeFile, TypeFluentd} {
			if err := Setup("github.com/realglobe-Inc", typ, lv, testOpt); err != nil {
				t.Fatal(err)
			}
		}
	}
}
