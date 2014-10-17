package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net"
	"net/http"
	"strconv"
)

func FreePort() (port int, err error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, erro.Wrap(err)
	}
	lis.Close()

	_, portStr, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		return 0, erro.Wrap(err)
	}

	port, err = strconv.Atoi(portStr)
	if err != nil {
		return 0, erro.Wrap(err)
	}

	return port, nil
}

// テスト用の HTTP サーバー。
type TestHttpServer struct {
	lis    net.Listener
	respCh chan *testHttpServerResponse
}

type testHttpServerResponse struct {
	status int
	header http.Header
	body   []byte
}

func NewTestHttpServer(port, chCap int) (*TestHttpServer, error) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	respCh := make(chan *testHttpServerResponse, chCap)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := <-respCh
		for key, values := range resp.header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.status)
		if resp.body != nil {
			w.Write(resp.body)
		}
	})

	go func() {
		http.Serve(lis, mux)
	}()

	return &TestHttpServer{lis, respCh}, nil
}

// 次に返させるレスポンスを与える。
func (server *TestHttpServer) AddResponse(status int, header http.Header, body []byte) {
	server.respCh <- &testHttpServerResponse{status, header, body}
}

// 待ち受けソケットを閉じる。
func (server *TestHttpServer) Close() error {
	if err := server.lis.Close(); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
