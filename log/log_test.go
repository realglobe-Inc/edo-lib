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
