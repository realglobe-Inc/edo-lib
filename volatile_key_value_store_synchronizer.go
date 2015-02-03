package driver

import (
	"reflect"
	"time"
)

type synchronizedVolatileKeyValueStore synchronizedDriver

type volatilePutRequest struct {
	key      string
	val      interface{}
	expiDate time.Time

	newCaStmpCh chan *Stamp
}

type entryRequest struct {
	eKey string

	eValCh chan string
}

type setEntryRequest struct {
	eKey      string
	eVal      string
	eExpiDate time.Time
}

type getAndSetEntryRequest struct {
	key       string
	caStmp    *Stamp
	eKey      string
	eVal      string
	eExpiDate time.Time

	valCh       chan interface{}
	newCaStmpCh chan *Stamp
}

type putIfEnteredRequest struct {
	key      string
	val      interface{}
	expiDate time.Time
	eKey     string
	eVal     string

	enteredCh   chan bool
	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedVolatileKeyValueStore(base ConcurrentVolatileKeyValueStore) *synchronizedVolatileKeyValueStore {
	return (*synchronizedVolatileKeyValueStore)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&kvsGetRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*kvsGetRequest)
			val, newCaStmp, err := base.Get(req.key, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.valCh <- val
				req.newCaStmpCh <- newCaStmp
			}
		},
		reflect.TypeOf(&volatilePutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*volatilePutRequest)
			newCaStmp, err := base.Put(req.key, req.val, req.expiDate)
			if err != nil {
				errCh <- err
			} else {
				req.newCaStmpCh <- newCaStmp
			}
		},
		reflect.TypeOf(&removeRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*removeRequest)
			errCh <- base.Remove(req.key)
		},
		reflect.TypeOf(&entryRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*entryRequest)
			eVal, err := base.Entry(req.eKey)
			if err != nil {
				errCh <- err
			} else {
				req.eValCh <- eVal
			}
		},
		reflect.TypeOf(&setEntryRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*setEntryRequest)
			errCh <- base.SetEntry(req.eKey, req.eVal, req.eExpiDate)
		},
		reflect.TypeOf(&getAndSetEntryRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*getAndSetEntryRequest)
			val, newCaStmp, err := base.GetAndSetEntry(req.key, req.caStmp, req.eKey, req.eVal, req.eExpiDate)
			if err != nil {
				errCh <- err
			} else {
				req.valCh <- val
				req.newCaStmpCh <- newCaStmp
			}
		},
		reflect.TypeOf(&putIfEnteredRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*putIfEnteredRequest)
			entered, newCaStmp, err := base.PutIfEntered(req.key, req.val, req.expiDate, req.eKey, req.eVal)
			if err != nil {
				errCh <- err
			} else {
				req.enteredCh <- entered
				req.newCaStmpCh <- newCaStmp
			}
		},
	}))
}

func (drv *synchronizedVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	valCh := make(chan interface{}, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&kvsGetRequest{key, caStmp, valCh, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return <-valCh, newCaStmp, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (drv *synchronizedVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&volatilePutRequest{key, val, expiDate, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return newCaStmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (drv *synchronizedVolatileKeyValueStore) Remove(key string) error {
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}

func (drv *synchronizedVolatileKeyValueStore) Entry(eKey string) (eVal string, err error) {
	eValCh := make(chan string, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&entryRequest{eKey, eValCh}, errCh}
	select {
	case eVal := <-eValCh:
		return eVal, nil
	case err := <-errCh:
		return "", err
	}
}

func (drv *synchronizedVolatileKeyValueStore) SetEntry(eKey, eVal string, eExpiDate time.Time) error {
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&setEntryRequest{eKey, eVal, eExpiDate}, errCh}
	return <-errCh
}

func (drv *synchronizedVolatileKeyValueStore) GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error) {
	valCh := make(chan interface{}, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&getAndSetEntryRequest{key, caStmp, eKey, eVal, eExpiDate, valCh, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return <-valCh, newCaStmp, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (drv *synchronizedVolatileKeyValueStore) PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error) {
	enteredCh := make(chan bool, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&putIfEnteredRequest{key, val, expiDate, eKey, eVal, enteredCh, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return <-enteredCh, newCaStmp, nil
	case err := <-errCh:
		return false, nil, err
	}
}
