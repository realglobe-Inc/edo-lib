package driver

import (
	"reflect"
)

// スレッドセーフにする。

// 非キャッシュ用。
func newSynchronizedIdProviderLister(reg IdProviderLister) IdProviderLister {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedIdProvidersRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedIdProvidersRequest)
			idps, err := reg.IdProviders()
			if err != nil {
				errCh <- err
			} else {
				req.idpsCh <- idps
			}
		},
	})
}

type synchronizedIdProvidersRequest struct {
	idpsCh chan []*IdProvider
}

func (reg *synchronizedDriver) IdProviders() ([]*IdProvider, error) {
	idpsCh := make(chan []*IdProvider, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedIdProvidersRequest{idpsCh}, errCh}
	select {
	case idps := <-idpsCh:
		return idps, nil
	case err := <-errCh:
		return nil, err
	}
}

// キャッシュ用。
func newSynchronizedDatedIdProviderLister(reg DatedIdProviderLister) DatedIdProviderLister {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedIdProvidersRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedIdProvidersRequest)
			idps, stmp, err := reg.StampedIdProviders(req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.idpsCh <- idps
				req.stmpCh <- stmp
			}
		},
	})
}

type synchronizedStampedIdProvidersRequest struct {
	caStmp *Stamp

	idpsCh chan []*IdProvider
	stmpCh chan *Stamp
}

func (reg *synchronizedDriver) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	idpsCh := make(chan []*IdProvider, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedIdProvidersRequest{caStmp, idpsCh, stmpCh}, errCh}
	select {
	case idps := <-idpsCh:
		return idps, <-stmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}
