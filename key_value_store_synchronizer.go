package driver

import (
	"reflect"
)

// スレッドセーフにする。

// 非キャッシュ用。
func newSynchronizedKeyValueStore(reg keyValueStore) keyValueStore {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedGetRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedGetRequest)
			if value, err := reg.get(req.key); err != nil {
				errCh <- err
			} else {
				req.valueCh <- value
			}
		},
		reflect.TypeOf(&synchronizedPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedPutRequest)
			errCh <- reg.put(req.key, req.value)
		},
		reflect.TypeOf(&synchronizedRemoveRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedRemoveRequest)
			errCh <- reg.remove(req.key)
		},
	})
}

type synchronizedGetRequest struct {
	key     string
	valueCh chan interface{}
}

type synchronizedPutRequest struct {
	key   string
	value interface{}
}

type synchronizedRemoveRequest struct {
	key string
}

func (reg *synchronizedDriver) get(key string) (value interface{}, err error) {
	valueCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedGetRequest{key, valueCh}, errCh}
	select {
	case value := <-valueCh:
		return value, nil
	case err := <-errCh:
		return nil, err
	}
}

func (reg *synchronizedDriver) put(key string, value interface{}) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedPutRequest{key, value}, errCh}
	return <-errCh
}

func (reg *synchronizedDriver) remove(key string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveRequest{key}, errCh}
	return <-errCh
}

// キャッシュ用。
func newSynchronizedDatedKeyValueStore(reg datedKeyValueStore) datedKeyValueStore {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedGetRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedGetRequest)
			value, stmp, err := reg.stampedGet(req.key, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.valueCh <- value
				req.stmpCh <- stmp
			}
		},
		reflect.TypeOf(&synchronizedStampedPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedPutRequest)
			stmp, err := reg.stampedPut(req.key, req.value)
			if err != nil {
				errCh <- err
			} else {
				req.stmpCh <- stmp
			}
		},
		reflect.TypeOf(&synchronizedRemoveRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedRemoveRequest)
			errCh <- reg.remove(req.key)
		},
	})
}

type synchronizedStampedGetRequest struct {
	key    string
	caStmp *Stamp

	valueCh chan interface{}
	stmpCh  chan *Stamp
}

type synchronizedStampedPutRequest struct {
	key   string
	value interface{}

	stmpCh chan *Stamp
}

func (reg *synchronizedDriver) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	valueCh := make(chan interface{}, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedGetRequest{key, caStmp, valueCh, stmpCh}, errCh}
	select {
	case value := <-valueCh:
		return value, <-stmpCh, nil
	case err := <-errCh:
		return "", nil, err
	}
}

func (reg *synchronizedDriver) stampedPut(key string, value interface{}) (newCaStmp *Stamp, err error) {
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedPutRequest{key, value, stmpCh}, errCh}
	select {
	case stmp := <-stmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}
