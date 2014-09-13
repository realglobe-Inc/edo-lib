package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"path/filepath"
	"time"
)

// 非キャッシュ用。
func NewFileJobRegistry(path string) JobRegistry {
	return newFileDriver(path)
}

func (reg *fileDriver) Result(jobId string) (*JobResult, error) {
	path := filepath.Join(reg.path, jobId+".json")

	var res JobResult
	if err := readFromJson(path, &res); err != nil {
		return nil, erro.Wrap(err)
	}

	if res.Status == 0 {
		return nil, nil
	}
	return &res, nil
}
func (reg *fileDriver) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	path := filepath.Join(reg.path, jobId+".json")

	if err := writeToJson(path, res); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
