package driver

import (
	"reflect"
)

type synchronizedListedRawDataStore struct {
	*synchronizedDriver
	base ListedRawDataStore
}

type keysRequest struct {
	caStmp *Stamp

	keysCh      chan map[string]bool
	newCaStmpCh chan *Stamp
}

type getRequest struct {
	key    string
	caStmp *Stamp

	dataCh      chan []byte
	newCaStmpCh chan *Stamp
}

type putRequest struct {
	key  string
	data []byte

	newCaStmpCh chan *Stamp
}

type removeRequest struct {
	key string
}

// もちろん、スレッドセーフ。
func newSynchronizedListedRawDataStore(base ListedRawDataStore) *synchronizedListedRawDataStore {
	return &synchronizedListedRawDataStore{newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&keysRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*keysRequest)
			keys, newCaStmp, err := base.Keys(req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.keysCh <- keys
				req.newCaStmpCh <- newCaStmp
			}
		},
		reflect.TypeOf(&getRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*getRequest)
			data, newCaStmp, err := base.Get(req.key, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.dataCh <- data
				req.newCaStmpCh <- newCaStmp
			}
		},
		reflect.TypeOf(&putRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*putRequest)
			newCaStmp, err := base.Put(req.key, req.data)
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
	}), base}
}

func (drv *synchronizedListedRawDataStore) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
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

func (drv *synchronizedListedRawDataStore) Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error) {
	dataCh := make(chan []byte, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&getRequest{key, caStmp, dataCh, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return <-dataCh, newCaStmp, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (drv *synchronizedListedRawDataStore) Put(key string, data []byte) (*Stamp, error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&putRequest{key, data, newCaStmpCh}, errCh}
	select {
	case newCaStmp := <-newCaStmpCh:
		return newCaStmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (drv *synchronizedListedRawDataStore) Remove(key string) error {
	errCh := make(chan error, 1)
	drv.reqCh <- &synchronizedRequest{&removeRequest{key}, errCh}
	return <-errCh
}

func (drv *synchronizedListedRawDataStore) Close() error {
	drv.synchronizedDriver.close()
	return drv.base.Close()
}
