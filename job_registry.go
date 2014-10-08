package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// ジョブ管理。
type JobRegistry interface {
	// 実行結果を取得する。
	Result(jobId string, caStmp *Stamp) (jobRes *JobResult, newCaStmp *Stamp, err error)

	// 実行結果を登録する。
	AddResult(jobId string, jobRes *JobResult, expiDate time.Time) (newCaStmp *Stamp, err error)
}

type JobResult struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// 骨組み。
// バックエンドのジョブ結果ごとに保存。
type jobRegistry struct {
	base TimeLimitedKeyValueStore
}

func newJobRegistry(base TimeLimitedKeyValueStore) *jobRegistry {
	return &jobRegistry{base}
}

func (reg *jobRegistry) Result(jobId string, caStmp *Stamp) (jobRes *JobResult, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(jobId, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.(*JobResult), newCaStmp, nil
}

func (reg *jobRegistry) AddResult(jobId string, jobRes *JobResult, expiDate time.Time) (newCaStmp *Stamp, err error) {
	return reg.base.Put(jobId, jobRes, expiDate)
}
