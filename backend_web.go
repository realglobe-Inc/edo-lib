package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"path"
	"time"
)

// JavaScript.
func NewWebJsBackendRegistry(addr string, ssl bool) (JsBackendRegistry, error) {
	return newWebDriver(addr, ssl)
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
		date := resp.Header.Get("Date")
		if date == "" {
			date = resp.Header.Get("Last-Modified")
		}
		expiDate := resp.Header.Get("Expires")
		etag := resp.Header.Get("ETag")

		newCaStmp := &Stamp{Digest: etag}
		if date == "" {
			newCaStmp.Date = time.Now()
		} else {
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

// ID プロバイダ.
func NewWebIdProviderBackend(addr string, ssl bool) (IdProviderBackend, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	req, err := http.NewRequest("GET", reg.prefix, nil)
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
		date := resp.Header.Get("Date")
		if date == "" {
			date = resp.Header.Get("Last-Modified")
		}
		expiDate := resp.Header.Get("Expires")
		etag := resp.Header.Get("ETag")

		stmp := &Stamp{Digest: etag}
		if date == "" {
			stmp.Date = time.Now()
		} else {
			stmp.Date, err = time.Parse(time.RFC1123, date)
			if err != nil {
				return nil, nil, erro.Wrap(err)
			}
		}
		if expiDate != "" {
			stmp.ExpiDate, err = time.Parse(time.RFC1123, expiDate)
			if err != nil {
				return nil, nil, erro.Wrap(err)
			}
		}

		switch resp.StatusCode {
		case http.StatusNotModified:
			return nil, stmp, nil
		case http.StatusOK:
			var idps []*IdProvider
			if err := json.NewDecoder(resp.Body).Decode(&idps); err != nil {
				return nil, nil, erro.Wrap(err)
			}
			return idps, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return nil, nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
