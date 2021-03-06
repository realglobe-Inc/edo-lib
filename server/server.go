// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// HTTP サーバー関係。
package server

import (
	"math/rand"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/realglobe-Inc/go-lib/erro"
)

// SocketType が tcp のときに追加で必要な関数。
type TcpParameter interface {
	// tcp のポート番号。
	SocketPort() int
}

// SocketType が unix のときに追加で必要な関数。
type UnixParameter interface {
	// unix のファイルパス。
	SocketPath() string
}

// デバッグに使える関数。
type DebugParameter interface {
	ShutdownChannel() chan struct{}
}

// 冷却期間の最大値。
var MaxSleepTime = time.Minute

// この期間以上エラー無しで動作したら冷却期間を 0 にする。
var ResetInterval = time.Minute

// サーバーを立てる。
func Serve(handler http.Handler, socType, protType string, param interface{}) error {
	// 冷却期間の揺らぎの最大値。
	// 冷却期間は前回の 2 倍に一様乱数の揺らぎを加えたものになる。
	const fluct = 100 * time.Millisecond

	var serv func(net.Listener, http.Handler) error
	switch protType {
	case "http":
		serv = http.Serve
	case "fcgi":
		serv = fcgi.Serve
	default:
		return erro.New("unsupported protocol type " + protType)
	}

	var shutCh chan struct{}
	if p, ok := param.(DebugParameter); ok {
		shutCh = p.ShutdownChannel()
	}
	if shutCh == nil {
		shutCh = make(chan struct{}, 5) // 余裕を持って。
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
			log.Info("Signal ", sig, " was detected")
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

	var sleepTime time.Duration = 0
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
				p, ok := param.(UnixParameter)
				if !ok {
					return false, erro.New("SocketPath function is not implemented")
				}
				l, err = net.Listen("unix", p.SocketPath())
				if err != nil {
					return true, erro.Wrap(err)
				} else if err := os.Chmod(p.SocketPath(), 0777); err != nil {
					return true, erro.Wrap(err)
				}
				log.Info("Wait on UNIX socket " + p.SocketPath())
			case "tcp":
				p, ok := param.(TcpParameter)
				if !ok {
					return false, erro.New("SocketPort function is not implemented")
				}
				l, err = net.Listen("tcp", ":"+strconv.Itoa(p.SocketPort()))
				if err != nil {
					return true, erro.Wrap(err)
				}
				log.Info("Wait on TCP socket ", p.SocketPort())
			default:
				return false, erro.New("unsupported socket type " + socType)
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
				log.Debug("Service starts")
				defer log.Debug("Service exits")
				return erro.Wrap(serv(l, handler))
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
				if end.Sub(start) > ResetInterval {
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

			sleepTime = nextSleepTime(sleepTime, fluct, MaxSleepTime)
			log.Info("Sleep ", sleepTime)

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

func nextSleepTime(cur, fluc, max time.Duration) time.Duration {
	next := 2*cur + time.Duration(rand.Int63n(int64(fluc)))
	if next > max {
		next = max
	}
	return next
}
