package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
)

// スレッドセーフでないデータ読み書きドライバーをスレッドセーフにする。
// 方法は非並列化。
type synchronizedDriver struct {
	reqCh chan *synchronizedRequest
}

type synchronizedRequest struct {
	req   interface{}
	errCh chan<- error
}

// キューの容量。
// 別に 0 でも良いが、そうすると毎回切り替えが必要になる。
const defChCap = 100

func newSynchronizedDriver(hndls map[reflect.Type]func(interface{}, chan<- error)) *synchronizedDriver {
	drv := &synchronizedDriver{
		make(chan *synchronizedRequest, defChCap),
	}

	go func() {
		for {
			drv.serve(hndls)
		}
	}()

	return drv
}

func (drv *synchronizedDriver) serve(hndls map[reflect.Type]func(interface{}, chan<- error)) {
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

	req := <-drv.reqCh
	errCh = req.errCh
	hndl := hndls[reflect.TypeOf(req.req)]
	if hndl != nil {
		hndl(req.req, errCh)
	}
}
