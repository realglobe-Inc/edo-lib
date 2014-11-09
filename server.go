package util

import (
	"errors"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"math/rand"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

var invalidProtocol = errors.New("invalid protocol.")

func Serve(socType, socPath string, socPort int, protType string, routes map[string]HandlerFunc) error {

	shutCh := make(chan struct{}, 1)

	// SIGINT、SIGKILL、SIGTERM を受け取ったら終了。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Info("Signal ", sig, " was caught.")
		shutCh <- struct{}{}
	}()

	// エラー発生時の冷却期間はだんだん延ばす。
	var sleepTime time.Duration = 0
	// 一定期間以上エラー無しで動作したら冷却期間をリセットする。
	resetInterval := time.Minute
	for {
		if brk, err := func() (brk bool, err error) {

			var lis net.Listener
			defer func() {
				if lis != nil {
					lis.Close()
				}
			}()

			switch socType {
			case "unix":
				lis, err = net.Listen("unix", socPath)
				if err != nil {
					return false, erro.Wrap(err)
				}
				if err := os.Chmod(socPath, 0777); err != nil {
					return false, erro.Wrap(err)
				}
				log.Info("Wait on UNIX socket " + socPath + ".")
			case "tcp":
				lis, err = net.Listen("tcp", ":"+strconv.Itoa(socPort))
				if err != nil {
					return false, erro.Wrap(err)
				}
				log.Info("Wait on TCP socket ", socPort, ".")
			default:
				return true, erro.New("invalid socket type " + socType + ".")
			}

			stopCh := make(chan struct{}, 1)
			subShutCh := make(chan bool, 1)
			go func() {
				select {
				case <-shutCh:
					shutCh <- struct{}{}
					subShutCh <- true
					lis.Close()
				case <-stopCh:
					subShutCh <- false
				}
			}()
			defer func() { stopCh <- struct{}{} }()

			start := time.Now()
			if err := serveCore(protType, routes, lis); err != nil {
				err := erro.Wrap(err)

				// 正常な終了処理としてソケットが閉じられたかもしれないので調べる。
				select {
				case <-subShutCh:
					return true, nil
				default:
				}

				stopCh <- struct{}{}
				brk = <-subShutCh

				if brk || erro.Unwrap(err) == invalidProtocol {
					// どうしようもない。
					return true, err
				}

				end := time.Now()
				if end.Sub(start) > resetInterval {
					sleepTime = 0
				}

				return false, err
			}

			stopCh <- struct{}{}
			return <-subShutCh, nil
		}(); brk {
			return erro.Wrap(err)
		} else {
			if err != nil {
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}

			sleepTime = serverNextSleepTime(sleepTime, resetInterval)
			log.Info("Retry after ", sleepTime)

			timeCh := time.After(sleepTime)
			select {
			case <-shutCh:
				return nil
			case <-timeCh:
			}
		}
	}

	return nil
}

func serverNextSleepTime(cur, max time.Duration) time.Duration {
	next := 2*cur + time.Duration(rand.Int63n(int64(time.Second)))
	if next >= max {
		next = time.Minute
	}
	return next
}

func serveCore(protType string, routes map[string]HandlerFunc, lis net.Listener) error {
	var serv func(net.Listener, http.Handler) error
	switch protType {
	case "http":
		serv = http.Serve
	case "fcgi":
		serv = fcgi.Serve
	default:
		return erro.Wrap(invalidProtocol)
	}

	log.Debug("Server starts.")
	defer log.Debug("Server exits.")

	mux := http.NewServeMux()

	for path, handler := range routes {
		mux.HandleFunc(path, serverPanicErrorWrapper(handler))
	}

	if err := serv(lis, mux); err != nil {
		return erro.Wrap(err)
	}

	return nil
}

// パニックとエラーの処理をまとめる。
func serverPanicErrorWrapper(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// panic時にプロセス終了しないようにrecoverする
		defer func() {
			if rcv := recover(); rcv != nil {
				buff := make([]byte, 8192)
				stackLen := runtime.Stack(buff, false)
				stack := string(buff[:stackLen])
				err := erro.Wrap(NewPanicWrapper(rcv, stack))

				log.Err(erro.Unwrap(err))
				log.Debug(err)

				body := ErrorToResponseJson(err)
				w.Header().Set("Content-Type", ContentTypeJson)
				w.Header().Set("Content-Length", strconv.Itoa(len(body)))
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(body)
				return
			}
		}()

		//////////////////////////////
		LogRequest(r, true)
		//////////////////////////////

		if err := handler(w, r); err != nil {
			err = erro.Wrap(err)
			log.Err(erro.Unwrap(err))
			log.Debug(err)

			var status int
			switch e := erro.Unwrap(err).(type) {
			case *HttpStatusError:
				status = e.Status()
			default:
				status = http.StatusInternalServerError
			}
			body := ErrorToResponseJson(err)
			w.Header().Set("Content-Type", ContentTypeJson)
			w.Header().Set("Content-Length", strconv.Itoa(len(body)))
			w.WriteHeader(status)
			w.Write(body)
			return
		}
	}
}
