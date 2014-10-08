package driver

import (
	"reflect"
)

type synchronizedKeyValueStore synchronizedDriver

type kvsGetRequest struct {
	key    string
	caStmp *Stamp

	valueCh     chan interface{}
	newCaStmpCh chan *Stamp
}

type kvsPutRequest struct {
	key   string
	value interface{}

	newCaStmpCh chan *Stamp
}

// もちろん、スレッドセーフ。
func newSynchronizedKeyValueStore(base KeyValueStore) *synchronizedKeyValueStore {
	return (*synchronizedKeyValueStore)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
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
		reflect.TypeOf(&kvsPutRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*kvsPutRequest)
			stmp, err := base.Put(req.key, req.value)
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

func (reg *synchronizedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
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

func (reg *synchronizedKeyValueStore) Put(key string, value interface{}) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&kvsPutRequest{key, value, newCaStmpCh}, errCh}
	select {
	case stmp := <-newCaStmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (reg *synchronizedKeyValueStore) Remove(key string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}
