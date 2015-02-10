package server

import (
	"github.com/realglobe-Inc/edo/util"
	logutil "github.com/realglobe-Inc/edo/util/log"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func init() {
	logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
}

// サーバーを停止させられるかどうかの検査。
func TestServerShutdown(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
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

// サーバーが諦めないことの検査。
func TestServerRestart(t *testing.T) {
	// ////////////////////////////////
	// logutil.SetupConsole("github.com/realglobe-Inc", level.ALL)
	// defer logutil.SetupConsole("github.com/realglobe-Inc", level.OFF)
	// ////////////////////////////////

	port, err := util.FreePort()
	if err != nil {
		t.Fatal(err)
	}

	shutCh1 := make(chan struct{}, 10)
	go TerminableServe("tcp", "", port, "http", map[string]HandlerFunc{"/": func(http.ResponseWriter, *http.Request) error { return nil }}, shutCh1, PanicErrorWrapper)
	defer func() { shutCh1 <- struct{}{} }()

	// 第一サーバー起動待ち。
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

	shutCh2 := make(chan struct{}, 10)
	go TerminableServe("tcp", "", port, "http", map[string]HandlerFunc{"/": func(http.ResponseWriter, *http.Request) error { return nil }}, shutCh2, PanicErrorWrapper)
	defer func() { shutCh2 <- struct{}{} }()

	// 第二サーバー起動失敗待ち。
	time.Sleep(100 * time.Millisecond)

	shutCh1 <- struct{}{}
	// 第一サーバー終了待ち。
	// 第二サーバー起動成功待ち。
	time.Sleep(200 * time.Millisecond)

	req, err = http.NewRequest("GET", "http://localhost:"+strconv.Itoa(port)+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Connection", "close")
	resp, err = (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error(resp)
	}

	shutCh2 <- struct{}{}
	// 第二サーバー終了待ち。
	time.Sleep(10 * time.Millisecond)

	if resp, err := (&http.Client{}).Get("http://localhost:" + strconv.Itoa(port) + "/"); err == nil {
		resp.Body.Close()
		t.Fatal(err)
	}
}
