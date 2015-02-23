package test

import (
	"bytes"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
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
type HttpServer struct {
	lis    net.Listener
	respCh chan *httpResponse
}

type httpResponse struct {
	status int
	header http.Header
	body   []byte

	reqCh chan<- *http.Request
}

func NewHttpServer(port int) (*HttpServer, error) {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	respCh := make(chan *httpResponse, 1024)

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
		for key, vals := range resp.header {
			for _, val := range vals {
				w.Header().Add(key, val)
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

	// 起動待ち。
	respCh <- &httpResponse{http.StatusOK, nil, nil, make(chan *http.Request, 1)}
	for i := time.Nanosecond; i < time.Second; i *= 2 {
		req, err := http.NewRequest("GET", "http://"+lis.Addr().String()+"/", nil)
		if err != nil {
			lis.Close()
			return nil, erro.Wrap(err)
		}
		req.Header.Set("Connection", "close")
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			// ちょっと待って再挑戦。
			time.Sleep(i)
			continue
		}
		// ちゃんとつながったので終わり。
		resp.Body.Close()
		return &HttpServer{lis, respCh}, nil
	}
	// 時間切れ。
	lis.Close()
	return nil, erro.New("time out")
}

func (this *HttpServer) Address() string {
	return this.lis.Addr().String()
}

// 次に返させるレスポンスを与える。
func (this *HttpServer) AddResponse(status int, header http.Header, body []byte) <-chan *http.Request {
	reqCh := make(chan *http.Request, 1)
	this.respCh <- &httpResponse{status, header, body, reqCh}
	return reqCh
}

// 待ち受けソケットを閉じる。
func (this *HttpServer) Close() error {
	if err := this.lis.Close(); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
