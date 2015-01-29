package util

import (
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func init() {
	SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
}

// サーバーを停止させられるかどうかの検査。
func TestServerShutdown(t *testing.T) {
	// ////////////////////////////////
	// SetupConsoleLog("github.com/realglobe-Inc", level.ALL)
	// defer SetupConsoleLog("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := FreePort()
	if err != nil {
		t.Fatal(err)
	}
	shutCh := make(chan struct{}, 10)
	go TerminableServe("tcp", "", port, "http", map[string]HandlerFunc{"/": func(http.ResponseWriter, *http.Request) error { return nil }}, shutCh, PanicErrorWrapper)
	defer func() { shutCh <- struct{}{} }()

	// サーバー起動待ち。
	time.Sleep(10 * time.Millisecond)

	req, err := http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Connection", "close")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	}

	shutCh <- struct{}{}
	// サーバー終了待ち。
	time.Sleep(10 * time.Millisecond)

	if resp, err := (&http.Client{}).Get("http://localhost:" + strconv.Itoa(port) + "/"); err == nil {
		resp.Body.Close()
		t.Fatal(err)
	}
}
