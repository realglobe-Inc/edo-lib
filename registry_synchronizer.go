package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
	"runtime"
	"time"
)

// スレッドセーフでないレジストリをスレッドセーフにする。

const defChCap = 1000

type synchronizedRegistry struct {
	reqCh chan *synchronizedRequest
}

type synchronizedRequest struct {
	req   interface{}
	errCh chan<- error
}

func newSynchronizedRegistry(hndls map[reflect.Type]func(interface{}, chan<- error)) *synchronizedRegistry {
	reg := &synchronizedRegistry{
		make(chan *synchronizedRequest, defChCap),
	}

	go func() {
		for {
			reg.serve(hndls)
		}
	}()

	return reg
}

func (reg *synchronizedRegistry) serve(hndls map[reflect.Type]func(interface{}, chan<- error)) {
	var errCh chan<- error
	defer func() {
		if rcv := recover(); rcv != nil {
			buff := make([]byte, 8192)
			stackLen := runtime.Stack(buff, false)
			stack := string(buff[:stackLen])
			err := erro.Wrap(util.NewPanicWrapper(rcv, stack))

			if errCh != nil {
				errCh <- err
			} else {
				log.Err(erro.Unwrap(err))
				log.Debug(err)
			}
		}
	}()

	req := <-reg.reqCh
	errCh = req.errCh
	hndl := hndls[reflect.TypeOf(req.req)]
	if hndl != nil {
		hndl(req.req, errCh)
	}
}

// ログイン。
func NewSynchronizedLoginRegistry(reg LoginRegistry) LoginRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedUserRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedUserRequest)
			usrUuid, err := reg.User(req.accToken)
			if err != nil {
				errCh <- err
			} else {
				req.usrCh <- usrUuid
			}
		},
	})
}

type synchronizedUserRequest struct {
	accToken string

	usrCh chan string
}

func (reg *synchronizedRegistry) User(accToken string) (usrUuid string, err error) {
	usrCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedUserRequest{accToken, usrCh}, errCh}
	select {
	case usrUuid := <-usrCh:
		return usrUuid, nil
	case err := <-errCh:
		return "", err
	}
}

// JavaScript.
func NewSynchronizedJsRegistry(reg JsRegistry) JsRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
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

type synchronizedObjectRequest struct {
	dir     string
	objName string

	objCh chan *Object
}
type synchronizedAddObjectRequest struct {
	dir     string
	objName string
	obj     *Object
}
type synchronizedRemoveObjectRequest struct {
	dir     string
	objName string
}

func (reg *synchronizedRegistry) Object(dir, objName string) (*Object, error) {
	objCh := make(chan *Object, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedObjectRequest{dir, objName, objCh}, errCh}
	select {
	case obj := <-objCh:
		return obj, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedRegistry) AddObject(dir, objName string, obj *Object) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddObjectRequest{dir, objName, obj}, errCh}
	return <-errCh
}
func (reg *synchronizedRegistry) RemoveObject(dir, objName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveObjectRequest{dir, objName}, errCh}
	return <-errCh
}

// ユーザー情報。
func NewSynchronizedUserRegistry(reg UserRegistry) UserRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
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
func (reg *synchronizedRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
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
func (reg *synchronizedRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddAttributeRequest{usrUuid, attrName, attr}, errCh}
	return <-errCh
}
func (reg *synchronizedRegistry) RemoveAttribute(usrUuid, attrName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveAttributeRequest{usrUuid, attrName}, errCh}
	return <-errCh
}

// ジョブ。
func NewSynchronizedJobRegistry(reg JobRegistry) JobRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedResultRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedResultRequest)
			res, err := reg.Result(req.jobId)
			if err != nil {
				errCh <- err
			} else {
				req.resCh <- res
			}
		},
		reflect.TypeOf(&synchronizedAddResultRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddResultRequest)
			errCh <- reg.AddResult(req.jobId, req.res, req.deadline)
		},
	})
}

type synchronizedResultRequest struct {
	jobId string

	resCh chan *JobResult
}
type synchronizedAddResultRequest struct {
	jobId    string
	res      *JobResult
	deadline time.Time
}

func (reg *synchronizedRegistry) Result(jobId string) (*JobResult, error) {
	resCh := make(chan *JobResult, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedResultRequest{jobId, resCh}, errCh}
	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedRegistry) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddResultRequest{jobId, res, deadline}, errCh}
	return <-errCh
}

// 別名。
func NewSynchronizedNameRegistry(reg NameRegistry) NameRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedAddressRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddressRequest)
			addr, err := reg.Address(req.name)
			if err != nil {
				errCh <- err
			} else {
				req.addrCh <- addr
			}
		},
		reflect.TypeOf(&synchronizedAddressesRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedAddressesRequest)
			addrs, err := reg.Addresses(req.name)
			if err != nil {
				errCh <- err
			} else {
				req.addrsCh <- addrs
			}
		},
	})
}

type synchronizedAddressRequest struct {
	name string

	addrCh chan string
}
type synchronizedAddressesRequest struct {
	name string

	addrsCh chan []string
}

func (reg *synchronizedRegistry) Address(name string) (addr string, err error) {
	addrCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddressRequest{name, addrCh}, errCh}
	select {
	case addr := <-addrCh:
		return addr, nil
	case err := <-errCh:
		return "", err
	}
}

func (reg *synchronizedRegistry) Addresses(name string) (addrs []string, err error) {
	addrsCh := make(chan []string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddressesRequest{name, addrsCh}, errCh}
	select {
	case addrs := <-addrsCh:
		return addrs, nil
	case err := <-errCh:
		return nil, err
	}
}

// イベント。
func NewSynchronizedEventRegistry(reg EventRegistry) EventRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedRegistry) Handler(usrUuid, event string) (Handler, error) {
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
func (reg *synchronizedRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddHandlerRequest{usrUuid, event, hndl}, errCh}
	return <-errCh
}
func (reg *synchronizedRegistry) RemoveHandler(usrUuid, event string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedRemoveHandlerRequest{usrUuid, event}, errCh}
	return <-errCh
}

// サービス。
func NewSynchronizedServiceRegistry(reg ServiceRegistry) ServiceRegistry {
	return newSynchronizedRegistry(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&synchronizedServiceRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*synchronizedServiceRequest)
			servUuid, err := reg.Service(req.addr)
			if err != nil {
				errCh <- err
			} else {
				req.servCh <- servUuid
			}
		},
	})
}

type synchronizedServiceRequest struct {
	addr string

	servCh chan string
}

func (reg *synchronizedRegistry) Service(addr string) (servUuid string, err error) {
	servCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedServiceRequest{addr, servCh}, errCh}
	select {
	case servUuid := <-servCh:
		return servUuid, nil
	case err := <-errCh:
		return "", err
	}
}
