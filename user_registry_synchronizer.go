package driver

import (
	"reflect"
)

// 非キャッシュ用。
func NewSynchronizedUserRegistry(reg UserRegistry) UserRegistry {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedAttributesRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAttributesRequest)
			attrs, err := reg.Attributes(req.usrUuid)
			if err != nil {
				errCh <- err
			} else {
				req.attrsCh <- attrs
			}
		},
		reflect.TypeOf(&synchronizedAttributeRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAttributeRequest)
			attr, err := reg.Attribute(req.usrUuid, req.attrName)
			if err != nil {
				errCh <- err
			} else {
				req.attrCh <- attr
			}
		},
		reflect.TypeOf(&synchronizedAddAttributeRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddAttributeRequest)
			errCh <- reg.AddAttribute(req.usrUuid, req.attrName, req.attr)
		},
		reflect.TypeOf(&synchronizedRemoveAttributeRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedRemoveAttributeRequest)
			errCh <- reg.RemoveAttribute(req.usrUuid, req.attrName)
		},
	})
}

type synchronizedAttributesRequest struct {
	usrUuid string

	attrsCh chan map[string]interface{}
}
type synchronizedAttributeRequest struct {
	usrUuid  string
	attrName string

	attrCh chan interface{}
}
type synchronizedAddAttributeRequest struct {
	usrUuid  string
	attrName string
	attr     interface{}
}
type synchronizedRemoveAttributeRequest struct {
	usrUuid  string
	attrName string
}

func (reg *synchronizedDriver) Attributes(usrUuid string) (map[string]interface{}, error) {
	attrsCh := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAttributesRequest{usrUuid, attrsCh}, errCh}
	select {
	case attrs := <-attrsCh:
		return attrs, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedDriver) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAttributeRequest{usrUuid, attrName, attrCh}, errCh}
	select {
	case attr := <-attrCh:
		return attr, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedDriver) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddAttributeRequest{usrUuid, attrName, attr}, errCh}
	return <-errCh
}
func (reg *synchronizedDriver) RemoveAttribute(usrUuid, attrName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveAttributeRequest{usrUuid, attrName}, errCh}
	return <-errCh
}
