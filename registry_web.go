package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Web API を提供するレジストリから取ってくる。

type webDriver struct {
	prefix string
	*http.Client
}

func newWebDriver(addr string, ssl bool) (*webDriver, error) {
	var prefix string
	if ssl {
		prefix = "http://"
	} else {
		prefix = "https://"
	}
	prefix += addr

	client := &http.Client{}

	// 途中で死んだときの対処をしっかりやる方向で。
	// // 接続テスト。
	// resp, err := client.Head(prefix)
	// if err != nil {
	// 	return nil, erro.Wrap(err)
	// }
	// resp.Body.Close()

	return &webDriver{prefix, client}, nil
}

// ログイン。
func NewWebLoginRegistry(addr string, ssl bool) (LoginRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) User(accToken string) (usrUuid string, err error) {
	resp, err := reg.Get(reg.prefix + "/" + accToken)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&usrUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return usrUuid, nil
}

// JavaScript.
func NewWebJsRegistry(addr string, ssl bool) (JsRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) Object(dir, objName string) (*Object, error) {
	resp, err := reg.Get(reg.prefix + path.Join(dir, objName))
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
	var obj Object
	if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return nil, erro.Wrap(err)
	}
	return &obj, nil
}
func (reg *webDriver) AddObject(dir, objName string, obj *Object) error {
	buff, err := json.Marshal(obj)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+path.Join(dir, objName), bytes.NewReader(buff))
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
func (reg *webDriver) RemoveObject(dir, objName string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+path.Join(dir, objName), nil)
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

// ユーザー情報。
func NewWebUserRegistry(addr string, ssl bool) (UserRegistry, error) {
	return newWebDriver(addr, ssl)
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

// ジョブ。
func NewWebJobRegistry(addr string, ssl bool) (JobRegistry, error) {
	return newWebDriver(addr, ssl)
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

// 住所。
func NewWebNameRegistry(addr string, ssl bool) (NameRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) Address(name string) (addr string, err error) {
	resp, err := reg.Get(reg.prefix + "/node/" + name)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addr); err != nil {
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *webDriver) Addresses(name string) (addrs []string, err error) {
	resp, err := reg.Get(reg.prefix + "/tree/" + name)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addrs); err != nil {
		return nil, erro.Wrap(err)
	}
	return addrs, nil
}

// イベント。
func NewWebEventRegistry(addr string, ssl bool) (EventRegistry, error) {
	return newWebDriver(addr, ssl)
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

// サービス。
func NewWebServiceRegistry(addr string, ssl bool) (ServiceRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) Service(addr string) (servUuid string, err error) {
	resp, err := reg.Get(reg.prefix + "?end_point=" + url.QueryEscape(addr))
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&servUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return servUuid, nil
}

// ID プロバイダ。
func NewWebIdProviderLister(addr string, ssl bool) (IdProviderLister, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) IdProviders() ([]*IdProvider, error) {
	resp, err := reg.Get(reg.prefix)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var idps []*IdProvider
	if err := json.NewDecoder(resp.Body).Decode(&idps); err != nil {
		return nil, erro.Wrap(err)
	}
	return idps, nil
}
