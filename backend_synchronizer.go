package driver

import (
	"reflect"
)

// JavaScript.
func NewSynchronizedJsBackendRegistry(reg JsBackendRegistry) JsBackendRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedStampedObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedStampedObjectRequest)
			obj, stmp, err := reg.StampedObject(req.dir, req.objName, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.objCh <- obj
				req.stmpCh <- stmp
			}
		},
		reflect.TypeOf(&synchronizedObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedObjectRequest)
			obj, err := reg.Object(req.dir, req.objName)
			if err != nil {
				errCh <- err
			} else {
				req.objCh <- obj
			}
		},
		reflect.TypeOf(&synchronizedAddObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddObjectRequest)
			errCh <- reg.AddObject(req.dir, req.objName, req.obj)
		},
		reflect.TypeOf(&synchronizedRemoveObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedRemoveObjectRequest)
			errCh <- reg.RemoveObject(req.dir, req.objName)
		},
	})
}

type synchronizedStampedObjectRequest struct {
	dir     string
	objName string
	caStmp  *Stamp

	objCh  chan *Object
	stmpCh chan *Stamp
}

func (reg *synchronizedRegistry) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	objCh := make(chan *Object, 1)
	stmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedStampedObjectRequest{dir, objName, caStmp, objCh, stmpCh}, errCh}
	select {
	case obj := <-objCh:
		return obj, <-stmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

// ID プロバイダ。
func NewSynchronizedIdProviderBackend(reg IdProviderBackend) IdProviderBackend {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedRegistry) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
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
