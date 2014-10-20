package util

import (
	"bytes"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
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

	reqCh chan<- *http.Request
}

func NewTestHttpServer(port int) (*TestHttpServer, error) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	respCh := make(chan *testHttpServerResponse, 100)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ここを抜けるときに勝手に Close されるので、Close されても問題無いように置き換える。
		req := *r
		if buff, err := ioutil.ReadAll(r.Body); err != nil {
			err := erro.Wrap(err)
			log.Err(erro.Unwrap(err))
			log.Debug(err)
		} else {
			req.Body = ioutil.NopCloser(bytes.NewReader(buff))
		}

		resp := <-respCh
		resp.reqCh <- &req
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
func (server *TestHttpServer) AddResponse(status int, header http.Header, body []byte) <-chan *http.Request {
	reqCh := make(chan *http.Request, 1)
	server.respCh <- &testHttpServerResponse{status, header, body, reqCh}
	return reqCh
}

// 待ち受けソケットを閉じる。
func (server *TestHttpServer) Close() error {
	if err := server.lis.Close(); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
