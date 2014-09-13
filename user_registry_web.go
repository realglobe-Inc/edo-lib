package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebUserRegistry(prefix string) UserRegistry {
	return newWebDriver(prefix)
}

func (reg *webDriver) Attributes(usrUuid string) (map[string]interface{}, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var attrs map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&attrs); err != nil {
		return nil, erro.Wrap(err)
	}
	return attrs, nil
}
func (reg *webDriver) Attribute(usrUuid, attrName string) (interface{}, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid + "/" + attrName)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	//util.LogResponse(resp, true)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var attr interface{}
	if err := json.NewDecoder(resp.Body).Decode(&attr); err != nil {
		return nil, erro.Wrap(err)
	}
	return attr, nil
}
func (reg *webDriver) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	buff, err := json.Marshal(attr)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+usrUuid+"/"+attrName, bytes.NewReader(buff))
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
func (reg *webDriver) RemoveAttribute(usrUuid, attrName string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+"/"+usrUuid+"/"+attrName, nil)
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
