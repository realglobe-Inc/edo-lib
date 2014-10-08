package driver

import (
	"reflect"
	"time"
)

type synchronizedTimeLimitedKeyValueStore synchronizedDriver

type tlPutRequest struct {
	key      string
	value    interface{}
	expiDate time.Time

	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedTimeLimitedKeyValueStore(base TimeLimitedKeyValueStore) *synchronizedTimeLimitedKeyValueStore {
	return (*synchronizedTimeLimitedKeyValueStore)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&kvsGetRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*kvsGetRequest)
			value, stmp, err := base.Get(req.key, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.valueCh <- value
				req.newCaStmpCh <- stmp
			}
		},
		reflect.TypeOf(&tlPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*tlPutRequest)
			stmp, err := base.Put(req.key, req.value, req.expiDate)
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

func (reg *synchronizedTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	valueCh := make(chan interface{}, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&kvsGetRequest{key, caStmp, valueCh, newCaStmpCh}, errCh}
	select {
	case value := <-valueCh:
		return value, <-newCaStmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (reg *synchronizedTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&tlPutRequest{key, value, expiDate, newCaStmpCh}, errCh}
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
