package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"path"
	"time"
)

// 非キャッシュ用。
func NewWebJsRegistry(prefix string) JsRegistry {
	return newWebDriver(prefix)
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

// キャッシュ用。
func NewWebJsBackendRegistry(prefix string) JsBackendRegistry {
	return newWebDriver(prefix)
}

func (reg *webDriver) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	req, err := http.NewRequest("GET", reg.prefix+path.Join(dir, objName), nil)
	if caStmp != nil {
		req.Header.Set("If-None-Match", caStmp.Digest)
		req.Header.Set("If-Modified-Since", caStmp.Date.Format(time.RFC1123))
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil, nil
	case http.StatusNotModified, http.StatusOK:
		date := resp.Header.Get("Last-Modified")
		expiDate := resp.Header.Get("Expires")
		etag := resp.Header.Get("ETag")

		newCaStmp := &Stamp{Digest: etag}
		if date != "" {
			newCaStmp.Date, err = time.Parse(time.RFC1123, date)
			if err != nil {
				return nil, nil, erro.Wrap(err)
			}
		}
		if expiDate != "" {
			newCaStmp.ExpiDate, err = time.Parse(time.RFC1123, expiDate)
			if err != nil {
				return nil, nil, erro.Wrap(err)
			}
		}

		switch resp.StatusCode {
		case http.StatusNotModified:
			return nil, newCaStmp, nil
		case http.StatusOK:
			var obj Object
			if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
				return nil, nil, erro.Wrap(err)
			}
			return &obj, newCaStmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return nil, nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
