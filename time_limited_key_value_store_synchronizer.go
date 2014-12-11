package driver

import (
	"reflect"
	"time"
)

type synchronizedTimeLimitedKeyValueStore synchronizedDriver

type tlPutRequest struct {
	key      string
	val      interface{}
	expiDate time.Time

	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedTimeLimitedKeyValueStore(base TimeLimitedKeyValueStore) *synchronizedTimeLimitedKeyValueStore {
	return (*synchronizedTimeLimitedKeyValueStore)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	valCh := make(chan interface{}, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&kvsGetRequest{key, caStmp, valCh, newCaStmpCh}, errCh}
	select {
	case val := <-valCh:
		return val, <-newCaStmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (reg *synchronizedTimeLimitedKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&tlPutRequest{key, val, expiDate, newCaStmpCh}, errCh}
	select {
	case stmp := <-newCaStmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (reg *synchronizedTimeLimitedKeyValueStore) Remove(key string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}
