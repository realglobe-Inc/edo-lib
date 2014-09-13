package driver

import (
	"reflect"
)

// 非キャッシュ用。
func NewSynchronizedEventRegistry(reg EventRegistry) EventRegistry {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedHandlerRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedHandlerRequest)
			hndl, err := reg.Handler(req.usrUuid, req.event)
			if err != nil {
				errCh <- err
			} else {
				req.hndlCh <- hndl
			}
		},
		reflect.TypeOf(&synchronizedAddHandlerRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddHandlerRequest)
			errCh <- reg.AddHandler(req.usrUuid, req.event, req.hndl)
		},
		reflect.TypeOf(&synchronizedRemoveHandlerRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedRemoveHandlerRequest)
			errCh <- reg.RemoveHandler(req.usrUuid, req.event)
		},
	})
}

type synchronizedHandlerRequest struct {
	usrUuid string
	event   string

	hndlCh chan Handler
}
type synchronizedAddHandlerRequest struct {
	usrUuid string
	event   string
	hndl    Handler
}
type synchronizedRemoveHandlerRequest struct {
	usrUuid string
	event   string
}

func (reg *synchronizedDriver) Handler(usrUuid, event string) (Handler, error) {
	hndlCh := make(chan Handler, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedHandlerRequest{usrUuid, event, hndlCh}, errCh}
	select {
	case hndl := <-hndlCh:
		return hndl, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedDriver) AddHandler(usrUuid, event string, hndl Handler) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddHandlerRequest{usrUuid, event, hndl}, errCh}
	return <-errCh
}
func (reg *synchronizedDriver) RemoveHandler(usrUuid, event string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveHandlerRequest{usrUuid, event}, errCh}
	return <-errCh
}
