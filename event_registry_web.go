package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebEventRegistry(prefix string) EventRegistry {
	return newWebDriver(prefix)
}

func (reg *webDriver) Handler(usrUuid, event string) (Handler, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid + event)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var hndl Handler
	if err := json.NewDecoder(resp.Body).Decode(&hndl); err != nil {
		return nil, erro.Wrap(err)
	}
	return hndl, nil
}
func (reg *webDriver) AddHandler(usrUuid, event string, hndl Handler) error {
	buff, err := json.Marshal(hndl)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+usrUuid+event, bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Set("Content-Type", util.ContentTypeJson)

	//util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	//util.LogResponse(resp, true)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	return nil
}
func (reg *webDriver) RemoveHandler(usrUuid, event string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+"/"+usrUuid+event, nil)
	if err != nil {
		return erro.Wrap(err)
	}
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	return nil
}
