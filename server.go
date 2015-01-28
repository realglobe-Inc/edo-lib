package util

import (
	"errors"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"math/rand"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

var invalidProtocol = errors.New("invalid protocol.")

func Serve(socType, socPath string, socPort int, protType string, routes map[string]HandlerFunc) error {
	return TerminableServe(socType, socPath, socPort, protType, routes, make(chan struct{}, 1))
}

func TerminableServe(socType, socPath string, socPort int, protType string, routes map[string]HandlerFunc, shutCh chan struct{}) error {

	var serv func(net.Listener, http.Handler) error
	switch protType {
	case "http":
		serv = http.Serve
	case "fcgi":
		serv = fcgi.Serve
	default:
		return erro.New("invalid protocol type " + protType)
	}

	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, serverPanicErrorWrapper(handler))
	}

	if shutCh == nil {
		shutCh = make(chan struct{}, 1)
	}
	var lis net.Listener
	var lisLock sync.Mutex

	// SIGINT、SIGTERM を受け取ったら正常終了。
	sigCh := make(chan os.Signal, 10)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 別スレッドで終了を監視する。
	go func() {
		select {
		case sig := <-sigCh:
			log.Info("Signal ", sig, " was caught")
		case <-shutCh:
			log.Info("Shutdown was operated")
		}
		shutCh <- struct{}{}
		lisLock.Lock()
		l := lis
		lisLock.Unlock()

		if l != nil {
			l.Close()
			log.Info("Socket was closed")
		} else {
			log.Info("No socket to close")
		}
	}()

	// エラー発生時の冷却期間はだんだん延ばす。
	var sleepTime time.Duration = 0
	// 一定期間以上エラー無しで動作したら冷却期間をリセットする。
	resetInterval := time.Minute
	for {
		if retry, err := func() (retry bool, err error) {

			var l net.Listener
			defer func() {
				if l != nil {
					l.Close()
				}
			}()

			switch socType {
			case "unix":
				l, err = net.Listen("unix", socPath)
				if err != nil {
					return true, erro.Wrap(err)
				}
				if err := os.Chmod(socPath, 0777); err != nil {
					return true, erro.Wrap(err)
				}
				log.Info("Wait on UNIX socket " + socPath)
			case "tcp":
				l, err = net.Listen("tcp", ":"+strconv.Itoa(socPort))
				if err != nil {
					return true, erro.Wrap(err)
				}
				log.Info("Wait on TCP socket ", socPort)
			default:
				return false, erro.New("invalid socket type " + socType)
			}

			lisLock.Lock()
			lis = l
			lisLock.Unlock()
			select {
			case <-shutCh:
				// 既に終了信号を受け取っていた。
				shutCh <- struct{}{}
				return false, nil
			default:
			}

			start := time.Now()
			if err := func() error {
				log.Debug("Service starts.")
				defer log.Debug("Service exits.")
				return erro.Wrap(serv(l, mux))
			}(); err != nil {
				err := erro.Wrap(err)

				select {
				case <-shutCh:
					// 既に終了信号を受け取っていた。
					shutCh <- struct{}{}
					return false, nil
				default:
				}

				end := time.Now()
				if end.Sub(start) > resetInterval {
					sleepTime = 0
				}

				return true, err
			}

			return true, nil
		}(); !retry {
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
				shutCh <- struct{}{}
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
		LogRequest(level.DEBUG, r, true)
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
