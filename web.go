package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
	"time"
)

// Web API を提供するレジストリから取ってくる。

func logRequest(r *http.Request) {
	buff, _ := httputil.DumpRequest(r, true)
	log.Debug("Request: ", string(buff))
}
func logResponse(r *http.Response) {
	buff, _ := httputil.DumpResponse(r, true)
	log.Debug("Response: ", string(buff))
}

type skeletalWebRegistry struct {
	prefix string
	*http.Client
}

func newSkeletalWebRegistry(addr string, ssl bool) (*skeletalWebRegistry, error) {
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

	return &skeletalWebRegistry{prefix, client}, nil
}

// JavaScript.
type webJsRegistry struct {
	*skeletalWebRegistry
}

func NewWebJsRegistry(addr string, ssl bool) (JsRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webJsRegistry{base}, nil
}

func (reg *webJsRegistry) Object(dir, objName string) (*Object, error) {
	resp, err := reg.Get(reg.prefix + path.Join(dir, objName))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	//logResponse(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	var obj Object
	if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return nil, erro.Wrap(err)
	}
	return &obj, nil
}
func (reg *webJsRegistry) AddObject(dir, objName string, obj *Object) error {
	buff, err := json.Marshal(obj)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+path.Join(dir, objName), bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Add("Content-Type", util.ContentTypeJson)

	//logRequest(req)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	//logResponse(resp)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}
func (reg *webJsRegistry) RemoveObject(dir, objName string) error {
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
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}

// ログイン。
type webLoginRegistry struct {
	*skeletalWebRegistry
}

func NewWebLoginRegistry(addr string, ssl bool) (LoginRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webLoginRegistry{base}, nil
}

func (reg *webLoginRegistry) User(accToken string) (usrUuid string, err error) {
	resp, err := reg.Get(reg.prefix + "/" + accToken)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, ".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&usrUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return usrUuid, nil
}

// ユーザー情報。
type webUserRegistry struct {
	*skeletalWebRegistry
}

func NewWebUserRegistry(addr string, ssl bool) (UserRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webUserRegistry{base}, nil
}

func (reg *webUserRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	var attrs map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&attrs); err != nil {
		return nil, erro.Wrap(err)
	}
	return attrs, nil
}
func (reg *webUserRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid + "/" + attrName)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	var attr interface{}
	if err := json.NewDecoder(resp.Body).Decode(&attr); err != nil {
		return nil, erro.Wrap(err)
	}
	return attr, nil
}
func (reg *webUserRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	buff, err := json.Marshal(attr)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+usrUuid+"/"+attrName, bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Add("Content-Type", util.ContentTypeJson)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}
func (reg *webUserRegistry) RemoveAttribute(usrUuid, attrName string) error {
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
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}

// ジョブ。
type webJobRegistry struct {
	*skeletalWebRegistry
}

func NewWebJobRegistry(addr string, ssl bool) (JobRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webJobRegistry{base}, nil
}

func (reg *webJobRegistry) Result(usrUuid string, jobId uint64) (interface{}, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid + "/" + strconv.FormatUint(jobId, 10))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	var attr interface{}
	if err := json.NewDecoder(resp.Body).Decode(&attr); err != nil {
		return nil, erro.Wrap(err)
	}
	return attr, nil
}

type ResultPack struct {
	Res     interface{} `json:"result"`
	Timelim time.Time   `json:"deadlineit"`
}

func (reg *webJobRegistry) AddResult(usrUuid string, jobId uint64, res interface{}, deadline time.Time) error {
	buff, err := json.Marshal(&ResultPack{res, deadline})
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+usrUuid+"/"+strconv.FormatUint(jobId, 10), bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Add("Content-Type", util.ContentTypeJson)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}

// 住所。
type webNameRegistry struct {
	*skeletalWebRegistry
}

func NewWebNameRegistry(addr string, ssl bool) (NameRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webNameRegistry{base}, nil
}

func (reg *webNameRegistry) Address(name string) (addr string, err error) {
	resp, err := reg.Get(reg.prefix + "/" + name + "/address")
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, ".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addr); err != nil {
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *webNameRegistry) Addresses(name string) (addrs []string, err error) {
	resp, err := reg.Get(reg.prefix + "/" + name)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addrs); err != nil {
		return nil, erro.Wrap(err)
	}
	return addrs, nil
}

// イベント。
type webEventRegistry struct {
	*skeletalWebRegistry
}

func NewWebEventRegistry(addr string, ssl bool) (EventRegistry, error) {
	base, err := newSkeletalWebRegistry(addr, ssl)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &webEventRegistry{base}, nil
}

func (reg *webEventRegistry) Handler(usrUuid, event string) (Handler, error) {
	resp, err := reg.Get(reg.prefix + "/" + usrUuid + "/" + event)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, ".")
	}
	var hndl Handler
	if err := json.NewDecoder(resp.Body).Decode(&hndl); err != nil {
		return nil, erro.Wrap(err)
	}
	return hndl, nil
}
func (reg *webEventRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	buff, err := json.Marshal(hndl)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+usrUuid+"/"+event, bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}
	req.Header.Add("Content-Type", util.ContentTypeJson)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}
func (reg *webEventRegistry) RemoveHandler(usrUuid, event string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+"/"+usrUuid+"/"+event, nil)
	if err != nil {
		return erro.Wrap(err)
	}
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return erro.New("invalid status ", resp.StatusCode, ".")
	}
	return nil
}
