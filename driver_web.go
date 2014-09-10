package driver

import (
	"net/http"
)

// HTTP/HTTPS で取ってくるドライバ。

type webDriver struct {
	prefix string // https://localhost:8000 とか。
	*http.Client
}

func newWebDriver(prefix string) *webDriver {
	// テストはしない。
	// 相手が一時的に落ちていても、復旧すれば動くように。
	return &webDriver{prefix, &http.Client{}}
}
