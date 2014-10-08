package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
	"runtime"
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
	reg := &synchronizedDriver{
		make(chan *synchronizedRequest, defChCap),
	}

	go func() {
		for {
			reg.serve(hndls)
		}
	}()

	return reg
}

func (reg *synchronizedDriver) serve(hndls map[reflect.Type]func(interface{}, chan<- error)) {
	var errCh chan<- error
	defer func() {
		if rcv := recover(); rcv != nil {
			buff := make([]byte, 8192)
			stackLen := runtime.Stack(buff, false)
			stack := string(buff[:stackLen])
			err := erro.Wrap(util.NewPanicWrapper(rcv, stack))

			if errCh != nil {
				errCh <- err
			} else {
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}
		}
	}()

	req := <-reg.reqCh
	errCh = req.errCh
	hndl := hndls[reflect.TypeOf(req.req)]
	if hndl != nil {
		hndl(req.req, errCh)
	}
}
