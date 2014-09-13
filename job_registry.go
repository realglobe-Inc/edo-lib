package driver

import (
	"time"
)

// ジョブ管理。
type JobRegistry interface {
	// 実行結果を取得する。
	Result(jobId string) (*JobResult, error)

	// 実行結果を登録する。
	AddResult(jobId string, res *JobResult, deadline time.Time) error
}

type JobResult struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}
