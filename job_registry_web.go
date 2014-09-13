package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"time"
)

// 非キャッシュ用。
func NewWebJobRegistry(prefix string) JobRegistry {
	return newWebDriver(prefix)
}

func (reg *webDriver) Result(jobId string) (*JobResult, error) {
	resp, err := reg.Get(reg.prefix + "/" + jobId)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified || resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var res JobResult
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &res, nil
}

type resultPack struct {
	*JobResult
	Deadline time.Time `json:"deadline"`
}

func (reg *webDriver) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	buff, err := json.Marshal(&resultPack{res, deadline})
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+jobId, bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Set("Content-Type", util.ContentTypeJson)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	return nil
}
