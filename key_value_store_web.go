package driver

import (
	"bytes"
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// GET http://example.com/{key}
// でボディに JSON で value が返り、
// PUT http://example.com/{key}
// のボディに JSON で value を入れると書き換えられ、
// DELETE http://example.com/{key}
// で消せると想定する。

// 非キャッシュ用。
func newWebKeyValueStore(prefix string) keyValueStore {
	return newWebDriver(prefix)
}

func (reg *webDriver) get(key string) (value interface{}, err error) {
	req, err := http.NewRequest("GET", reg.prefix+"/"+key, nil)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		return nil, erro.Wrap(err)
	}
	return value, nil
}

func (reg *webDriver) put(key string, value interface{}) error {
	buff, err := json.Marshal(value)
	if err != nil {
		return erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+key, bytes.NewReader(buff))
	if err != nil {
		return erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}

	return nil
}

func (reg *webDriver) remove(key string) error {
	req, err := http.NewRequest("DELETE", reg.prefix+"/"+key, nil)
	if err != nil {
		return erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	if resp.StatusCode != http.StatusOK {
		return erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}

	return nil
}

// キャッシュ用。
func newWebDatedKeyValueStore(prefix string) datedKeyValueStore {
	return newWebDriver(prefix)
}

func (reg *webDriver) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	req, err := http.NewRequest("GET", reg.prefix+"/"+key, nil)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	if caStmp != nil {
		WriteStampToRequestHeader(caStmp, req.Header)
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
		stmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}

		switch resp.StatusCode {
		case http.StatusNotModified:
			return nil, stmp, nil
		case http.StatusOK:
			if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
				return nil, nil, erro.Wrap(err)
			}
			return value, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return nil, nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}

func (reg *webDriver) stampedPut(key string, value interface{}) (*Stamp, error) {
	buff, err := json.Marshal(value)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	req, err := http.NewRequest("PUT", reg.prefix+"/"+key, bytes.NewReader(buff))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusNotModified, http.StatusOK:
		newCaStmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, erro.Wrap(err)
		}

		switch resp.StatusCode {
		case http.StatusOK, http.StatusCreated:
			return newCaStmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
