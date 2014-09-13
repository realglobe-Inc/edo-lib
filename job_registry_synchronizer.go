package driver

import (
	"reflect"
	"time"
)

// 非キャッシュ用。
func NewSynchronizedJobRegistry(reg JobRegistry) JobRegistry {
	return newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
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

func (reg *synchronizedDriver) Result(jobId string) (*JobResult, error) {
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
func (reg *synchronizedDriver) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&synchronizedAddResultRequest{jobId, res, deadline}, errCh}
	return <-errCh
}
