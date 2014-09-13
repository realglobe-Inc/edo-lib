package driver

import (
	"reflect"
)

// 非キャッシュ用。
func NewSynchronizedNameRegistry(reg NameRegistry) NameRegistry {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedDriver) Address(name string) (addr string, err error) {
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

func (reg *synchronizedDriver) Addresses(name string) (addrs []string, err error) {
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
