package driver

import (
	"reflect"
)

// 非キャッシュ用。
func NewSynchronizedLoginRegistry(reg LoginRegistry) LoginRegistry {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedDriver) User(accToken string) (usrUuid string, err error) {
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
