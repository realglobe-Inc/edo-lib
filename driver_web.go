package driver

import (
	"net/http"
)

// スレッドセーフ。
type webDriver struct {
	prefix string // https://localhost:8000 とか。
	Client *http.Client
}

func newWebDriver(prefix string) *webDriver {
	// テストはしない。
	// 相手が一時的に落ちていても、復旧すれば動くように。
	return &webDriver{prefix, &http.Client{}}
}
