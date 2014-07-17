package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"runtime"
	"time"
)

// スレッドセーフでないレジストリをスレッドセーフにする。

const defChCap = 1000

// JavaScript.
type synchronizedJsRegistry struct {
	reqCh chan interface{}
}

type synchronizedObjectRequest struct {
	dir     string
	objName string

	objCh chan *Object
	errCh chan error
}
type synchronizedAddObjectRequest struct {
	dir     string
	objName string
	obj     *Object

	errCh chan error
}
type synchronizedRemoveObjectRequest struct {
	dir     string
	objName string

	errCh chan error
}

func NewSynchronizedJsRegistry(base JsRegistry) JsRegistry {
	reg := &synchronizedJsRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedJsRegistry) Object(dir, objName string) (*Object, error) {
	objCh := make(chan *Object, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedObjectRequest{dir, objName, objCh, errCh}
	select {
	case obj := <-objCh:
		return obj, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedJsRegistry) AddObject(dir, objName string, obj *Object) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddObjectRequest{dir, objName, obj, errCh}
	return <-errCh
}
func (reg *synchronizedJsRegistry) RemoveObject(dir, objName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRemoveObjectRequest{dir, objName, errCh}
	return <-errCh
}

func (reg *synchronizedJsRegistry) serve(base JsRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedObjectRequest:
		errCh = req.errCh
		obj, err := base.Object(req.dir, req.objName)
		if err != nil {
			req.errCh <- err
		} else {
			req.objCh <- obj
		}
	case *synchronizedAddObjectRequest:
		errCh = req.errCh
		req.errCh <- base.AddObject(req.dir, req.objName, req.obj)
	case *synchronizedRemoveObjectRequest:
		errCh = req.errCh
		req.errCh <- base.RemoveObject(req.dir, req.objName)
	}
}

// ログイン。
type synchronizedLoginRegistry struct {
	reqCh chan interface{}
}

type synchronizedLoginRequest struct {
	accToken string

	usrCh chan string
	errCh chan error
}

func NewSynchronizedLoginRegistry(base LoginRegistry) LoginRegistry {
	reg := &synchronizedLoginRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedLoginRegistry) User(accToken string) (usrUuid string, err error) {
	usrCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedLoginRequest{accToken, usrCh, errCh}
	select {
	case usrUuid := <-usrCh:
		return usrUuid, nil
	case err := <-errCh:
		return "", err
	}
}

func (reg *synchronizedLoginRegistry) serve(base LoginRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedLoginRequest:
		errCh = req.errCh
		usrUuid, err := base.User(req.accToken)
		if err != nil {
			req.errCh <- err
		} else {
			req.usrCh <- usrUuid
		}
	}
}

// ユーザー情報。
type synchronizedUserRegistry struct {
	reqCh chan interface{}
}

type synchronizedAttributesRequest struct {
	usrUuid string

	attrsCh chan map[string]interface{}
	errCh   chan error
}
type synchronizedAttributeRequest struct {
	usrUuid  string
	attrName string

	attrCh chan interface{}
	errCh  chan error
}
type synchronizedAddAttributeRequest struct {
	usrUuid  string
	attrName string
	attr     interface{}

	errCh chan error
}
type synchronizedRemoveAttributeRequest struct {
	usrUuid  string
	attrName string

	errCh chan error
}

func NewSynchronizedUserRegistry(base UserRegistry) UserRegistry {
	reg := &synchronizedUserRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedUserRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
	attrsCh := make(chan map[string]interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAttributesRequest{usrUuid, attrsCh, errCh}
	select {
	case attrs := <-attrsCh:
		return attrs, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedUserRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAttributeRequest{usrUuid, attrName, attrCh, errCh}
	select {
	case attr := <-attrCh:
		return attr, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedUserRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddAttributeRequest{usrUuid, attrName, attr, errCh}
	return <-errCh
}
func (reg *synchronizedUserRegistry) RemoveAttribute(usrUuid, attrName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRemoveAttributeRequest{usrUuid, attrName, errCh}
	return <-errCh
}

func (reg *synchronizedUserRegistry) serve(base UserRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedAttributesRequest:
		errCh = req.errCh
		attrs, err := base.Attributes(req.usrUuid)
		if err != nil {
			req.errCh <- err
		} else {
			req.attrsCh <- attrs
		}
	case *synchronizedAttributeRequest:
		errCh = req.errCh
		attr, err := base.Attribute(req.usrUuid, req.attrName)
		if err != nil {
			req.errCh <- err
		} else {
			req.attrCh <- attr
		}
	case *synchronizedAddAttributeRequest:
		errCh = req.errCh
		req.errCh <- base.AddAttribute(req.usrUuid, req.attrName, req.attr)
	case *synchronizedRemoveAttributeRequest:
		errCh = req.errCh
		req.errCh <- base.RemoveAttribute(req.usrUuid, req.attrName)
	}
}

// ジョブ。
type synchronizedJobRegistry struct {
	reqCh chan interface{}
}

type synchronizedResultRequest struct {
	usrUuid string
	jobId   uint64

	resCh chan interface{}
	errCh chan error
}
type synchronizedAddResultRequest struct {
	usrUuid  string
	jobId    uint64
	res      interface{}
	deadline time.Time

	errCh chan error
}

func NewSynchronizedJobRegistry(base JobRegistry) JobRegistry {
	reg := &synchronizedJobRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedJobRegistry) Result(usrUuid string, jobId uint64) (interface{}, error) {
	resCh := make(chan interface{}, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedResultRequest{usrUuid, jobId, resCh, errCh}
	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedJobRegistry) AddResult(usrUuid string, jobId uint64, res interface{}, deadline time.Time) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddResultRequest{usrUuid, jobId, res, deadline, errCh}
	return <-errCh
}

func (reg *synchronizedJobRegistry) serve(base JobRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedResultRequest:
		errCh = req.errCh
		res, err := base.Result(req.usrUuid, req.jobId)
		if err != nil {
			req.errCh <- err
		} else {
			req.resCh <- res
		}
	case *synchronizedAddResultRequest:
		errCh = req.errCh
		req.errCh <- base.AddResult(req.usrUuid, req.jobId, req.res, req.deadline)
	}
}

// 別名。
type synchronizedNameRegistry struct {
	reqCh chan interface{}
}

type synchronizedAddressRequest struct {
	name string

	addrCh chan string
	errCh  chan error
}
type synchronizedAddressesRequest struct {
	name string

	addrsCh chan []string
	errCh   chan error
}

func NewSynchronizedNameRegistry(base NameRegistry) NameRegistry {
	reg := &synchronizedNameRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedNameRegistry) Address(name string) (addr string, err error) {
	addrCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddressRequest{name, addrCh, errCh}
	select {
	case addr := <-addrCh:
		return addr, nil
	case err := <-errCh:
		return "", err
	}
}

func (reg *synchronizedNameRegistry) Addresses(name string) (addrs []string, err error) {
	addrsCh := make(chan []string, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddressesRequest{name, addrsCh, errCh}
	select {
	case addrs := <-addrsCh:
		return addrs, nil
	case err := <-errCh:
		return nil, err
	}
}

func (reg *synchronizedNameRegistry) serve(base NameRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedAddressRequest:
		errCh = req.errCh
		addr, err := base.Address(req.name)
		if err != nil {
			req.errCh <- err
		} else {
			req.addrCh <- addr
		}
	case *synchronizedAddressesRequest:
		errCh = req.errCh
		addrs, err := base.Addresses(req.name)
		if err != nil {
			req.errCh <- err
		} else {
			req.addrsCh <- addrs
		}
	}
}

// イベント。
type synchronizedEventRegistry struct {
	reqCh chan interface{}
}

type synchronizedHandlerRequest struct {
	usrUuid string
	event   string

	hndlCh chan Handler
	errCh  chan error
}
type synchronizedAddHandlerRequest struct {
	usrUuid string
	event   string
	hndl    Handler

	errCh chan error
}
type synchronizedRemoveHandlerRequest struct {
	usrUuid string
	event   string

	errCh chan error
}

func NewSynchronizedEventRegistry(base EventRegistry) EventRegistry {
	reg := &synchronizedEventRegistry{
		make(chan interface{}, defChCap),
	}

	go func() {
		for {
			reg.serve(base)
		}
	}()

	return reg
}

func (reg *synchronizedEventRegistry) Handler(usrUuid, event string) (Handler, error) {
	hndlCh := make(chan Handler, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedHandlerRequest{usrUuid, event, hndlCh, errCh}
	select {
	case hndl := <-hndlCh:
		return hndl, nil
	case err := <-errCh:
		return nil, err
	}
}
func (reg *synchronizedEventRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedAddHandlerRequest{usrUuid, event, hndl, errCh}
	return <-errCh
}
func (reg *synchronizedEventRegistry) RemoveHandler(usrUuid, event string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRemoveHandlerRequest{usrUuid, event, errCh}
	return <-errCh
}

func (reg *synchronizedEventRegistry) serve(base EventRegistry) {
	var errCh chan error
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

	switch req := (<-reg.reqCh).(type) {
	case *synchronizedHandlerRequest:
		errCh = req.errCh
		hndl, err := base.Handler(req.usrUuid, req.event)
		if err != nil {
			req.errCh <- err
		} else {
			req.hndlCh <- hndl
		}
	case *synchronizedAddHandlerRequest:
		errCh = req.errCh
		req.errCh <- base.AddHandler(req.usrUuid, req.event, req.hndl)
	case *synchronizedRemoveHandlerRequest:
		errCh = req.errCh
		req.errCh <- base.RemoveHandler(req.usrUuid, req.event)
	}
}
