package driver

import (
	"reflect"
)

// スレッドセーフにする。

// 非キャッシュ用。
func NewSynchronizedServiceKeyRegistry(reg ServiceKeyRegistry) ServiceKeyRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedServiceKeyRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedServiceKeyRequest)
			key, err := reg.ServiceKey(req.servUuid)
			if err != nil {
				errCh <- err
			} else {
				req.keyCh <- key
			}
		},
	})
}

type synchronizedServiceKeyRequest struct {
	servUuid string
	keyCh    chan string
}

func (reg *synchronizedRegistry) ServiceKey(servUuid string) (key string, err error) {
	keyCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedServiceKeyRequest{servUuid, keyCh}, errCh}
	select {
	case key := <-keyCh:
		return key, nil
	case err := <-errCh:
		return "", err
	}
}

// キャッシュ用。
func NewSynchronizedDatedServiceKeyRegistry(reg DatedServiceKeyRegistry) DatedServiceKeyRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedServiceKeyRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedServiceKeyRequest)
			key, stmp, err := reg.StampedServiceKey(req.servUuid, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.keyCh <- key
				req.stmpCh <- stmp
			}
		},
	})
}

type synchronizedStampedServiceKeyRequest struct {
	servUuid string
	caStmp   *Stamp

	keyCh  chan string
	stmpCh chan *Stamp
}

func (reg *synchronizedRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error) {
	keyCh := make(chan string, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedServiceKeyRequest{servUuid, caStmp, keyCh, stmpCh}, errCh}
	select {
	case key := <-keyCh:
		return key, <-stmpCh, nil
	case err := <-errCh:
		return "", nil, err
	}
}
