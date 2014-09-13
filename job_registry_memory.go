package driver

import (
	"time"
)

// 非キャッシュ用。
type MemoryJobRegistry struct {
	ress map[string]*JobResult
}

func NewMemoryJobRegistry() *MemoryJobRegistry {
	return &MemoryJobRegistry{map[string]*JobResult{}}
}

func (reg *MemoryJobRegistry) Result(jobId string) (res *JobResult, err error) {
	return reg.ress[jobId], nil
}
func (reg *MemoryJobRegistry) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	reg.ress[jobId] = res
	return nil
}
