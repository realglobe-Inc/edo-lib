package driver

import (
	"reflect"
)

// スレッドセーフにする。

// 非キャッシュ用。
func NewSynchronizedServiceExplorer(reg ServiceExplorer) ServiceExplorer {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedServiceUuidRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedServiceUuidRequest)
			servUuid, err := reg.ServiceUuid(req.servUri)
			if err != nil {
				errCh <- err
			} else {
				req.servUuidCh <- servUuid
			}
		},
	})
}

type synchronizedServiceUuidRequest struct {
	servUri    string
	servUuidCh chan string
}

func (reg *synchronizedRegistry) ServiceUuid(servUri string) (servUuid string, err error) {
	servUuidCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedServiceUuidRequest{servUri, servUuidCh}, errCh}
	select {
	case servUuid := <-servUuidCh:
		return servUuid, nil
	case err := <-errCh:
		return "", err
	}
}

// キャッシュ用。
func NewSynchronizedDatedServiceExplorer(reg DatedServiceExplorer) DatedServiceExplorer {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedServiceUuidRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedServiceUuidRequest)
			servUuid, stmp, err := reg.StampedServiceUuid(req.servUri, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.servUuidCh <- servUuid
				req.stmpCh <- stmp
			}
		},
	})
}

type synchronizedStampedServiceUuidRequest struct {
	servUri string
	caStmp  *Stamp

	servUuidCh chan string
	stmpCh     chan *Stamp
}

func (reg *synchronizedRegistry) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	servUuidCh := make(chan string, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedServiceUuidRequest{servUri, caStmp, servUuidCh, stmpCh}, errCh}
	select {
	case servUuid := <-servUuidCh:
		return servUuid, <-stmpCh, nil
	case err := <-errCh:
		return "", nil, err
	}
}
