package driver

import (
	"reflect"
	"time"
)

type synchronizedVolatileKeyValueStore synchronizedDriver

type tlPutRequest struct {
	key      string
	val      interface{}
	expiDate time.Time

	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedVolatileKeyValueStore(base VolatileKeyValueStore) *synchronizedVolatileKeyValueStore {
	return (*synchronizedVolatileKeyValueStore)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&kvsGetRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*kvsGetRequest)
			val, stmp, err := base.Get(req.key, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.valCh <- val
				req.newCaStmpCh <- stmp
			}
		},
		reflect.TypeOf(&tlPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*tlPutRequest)
			stmp, err := base.Put(req.key, req.val, req.expiDate)
			if err != nil {
				errCh <- err
			} else {
				req.newCaStmpCh <- stmp
			}
		},
		reflect.TypeOf(&removeRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*removeRequest)
			errCh <- base.Remove(req.key)
		},
	}))
}

func (drv *synchronizedVolatileKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	valCh := make(chan interface{}, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&kvsGetRequest{key, caStmp, valCh, newCaStmpCh}, errCh}
	select {
	case val := <-valCh:
		return val, <-newCaStmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (drv *synchronizedVolatileKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&tlPutRequest{key, val, expiDate, newCaStmpCh}, errCh}
	select {
	case stmp := <-newCaStmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (drv *synchronizedVolatileKeyValueStore) Remove(key string) error {
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}
