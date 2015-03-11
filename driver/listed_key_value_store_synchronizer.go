package driver

import (
	"reflect"
)

type synchronizedListedKeyValueStore struct {
	*synchronizedDriver
	base ListedKeyValueStore
}

type kvsGetRequest struct {
	key    string
	caStmp *Stamp

	valCh       chan interface{}
	newCaStmpCh chan *Stamp
}

type kvsPutRequest struct {
	key string
	val interface{}

	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedListedKeyValueStore(base ListedKeyValueStore) *synchronizedListedKeyValueStore {
	return &synchronizedListedKeyValueStore{newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&keysRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*keysRequest)
			keys, stmp, err := base.Keys(req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.keysCh <- keys
				req.newCaStmpCh <- stmp
			}
		},
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
		reflect.TypeOf(&kvsPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*kvsPutRequest)
			stmp, err := base.Put(req.key, req.val)
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
	}), base}
}

func (drv *synchronizedListedKeyValueStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	keysCh := make(chan map[string]bool, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&keysRequest{caStmp, keysCh, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return <-keysCh, newCaStmp, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (drv *synchronizedListedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
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

func (drv *synchronizedListedKeyValueStore) Put(key string, val interface{}) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&kvsPutRequest{key, val, newCaStmpCh}, errCh}
	select {
	case stmp := <-newCaStmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (drv *synchronizedListedKeyValueStore) Remove(key string) error {
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}

func (drv *synchronizedListedKeyValueStore) Close() error {
	drv.synchronizedDriver.close()
	return drv.base.Close()
}
