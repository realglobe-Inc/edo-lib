package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"time"
)

// 非キャッシュ用。
func NewWebIdProviderLister(prefix string) IdProviderLister {
	return newWebDriver(prefix)
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

// キャッシュ用。
func NewWebDatedIdProviderLister(prefix string) DatedIdProviderLister {
	// TODO キャッシュの並列化。
	return newSynchronizedDatedIdProviderLister(newCachingDatedIdProviderLister(newWebDriver(prefix)))
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
		date := resp.Header.Get("Last-Modified")
		expiDate := resp.Header.Get("Expires")
		etag := resp.Header.Get("ETag")

		stmp := &Stamp{Digest: etag}
		if date != "" {
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
