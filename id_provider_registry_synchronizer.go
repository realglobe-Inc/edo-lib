package driver

import (
	"reflect"
)

// スレッドセーフにする。

// 非キャッシュ用。
func NewSynchronizedIdProviderRegistry(reg IdProviderRegistry) IdProviderRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedIdProviderQueryUriRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedIdProviderQueryUriRequest)
			queryUri, err := reg.IdProviderQueryUri(req.idpUuid)
			if err != nil {
				errCh <- err
			} else {
				req.queryUriCh <- queryUri
			}
		},
	})
}

type synchronizedIdProviderQueryUriRequest struct {
	idpUuid    string
	queryUriCh chan string
}

func (reg *synchronizedRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	queryUriCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedIdProviderQueryUriRequest{idpUuid, queryUriCh}, errCh}
	select {
	case queryUri := <-queryUriCh:
		return queryUri, nil
	case err := <-errCh:
		return "", err
	}
}

// キャッシュ用。
func NewSynchronizedDatedIdProviderRegistry(reg DatedIdProviderRegistry) DatedIdProviderRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedIdProviderQueryUriRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedIdProviderQueryUriRequest)
			queryUri, stmp, err := reg.StampedIdProviderQueryUri(req.idpUuid, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.queryUriCh <- queryUri
				req.stmpCh <- stmp
			}
		},
	})
}

type synchronizedStampedIdProviderQueryUriRequest struct {
	idpUuid string
	caStmp  *Stamp

	queryUriCh chan string
	stmpCh     chan *Stamp
}

func (reg *synchronizedRegistry) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	queryUriCh := make(chan string, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedIdProviderQueryUriRequest{idpUuid, caStmp, queryUriCh, stmpCh}, errCh}
	select {
	case queryUri := <-queryUriCh:
		return queryUri, <-stmpCh, nil
	case err := <-errCh:
		return "", nil, err
	}
}
