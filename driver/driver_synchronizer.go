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

package driver

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"reflect"
)

// スレッドセーフでないデータ読み書きドライバーをスレッドセーフにする。
// 方法は非並列化。
type synchronizedDriver struct {
	reqCh  chan *synchronizedRequest
	shutCh chan struct{}
}

type synchronizedRequest struct {
	req   interface{}
	errCh chan<- error
}

// キューの容量。
// 別に 0 でも良いが、そうすると毎回切り替えが必要になる。
const defChCap = 1024

func newSynchronizedDriver(hndls map[reflect.Type]func(interface{}, chan<- error)) *synchronizedDriver {
	drv := &synchronizedDriver{
		make(chan *synchronizedRequest, defChCap),
		make(chan struct{}, 1),
	}

	go func() {
		for drv.serve(hndls) {
		}
	}()

	return drv
}

// 返り値は終了信号を受け取ったときのみ false。
func (drv *synchronizedDriver) serve(hndls map[reflect.Type]func(interface{}, chan<- error)) (cont bool) {
	var errCh chan<- error
	defer func() {
		if rcv := recover(); rcv != nil {
			err := erro.New(rcv)

			if errCh != nil {
				errCh <- err
			} else {
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}
		}
	}()

	req, ok := <-drv.reqCh
	if !ok {
		drv.shutCh <- struct{}{}
		return false
	}

	errCh = req.errCh
	hndl := hndls[reflect.TypeOf(req.req)]
	if hndl != nil {
		hndl(req.req, errCh)
	}
	return true
}

func (drv *synchronizedDriver) close() {
	if drv.reqCh == nil {
		return
	}

	close(drv.reqCh)
	<-drv.shutCh
	// serve してるスレッドに reqCh が閉じてることが伝わった。

	drv.shutCh <- struct{}{}
	drv.reqCh = nil
}
