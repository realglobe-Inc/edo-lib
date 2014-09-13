package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"time"
)

// {
//   "id_providers": [
//     { "uuid": "aaaa-bbbb-cccc", "name": "リアルグローブ", "uri": "https://realglobe.jp/login" },
//     ...
//   ]
// }

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

	var res struct {
		Idps []*IdProvider `json:"id_providers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Idps, nil
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
		stmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}

		switch resp.StatusCode {
		case http.StatusNotModified:
			return nil, stmp, nil
		case http.StatusOK:
			var res struct {
				Idps []*IdProvider `json:"id_providers"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				return nil, nil, erro.Wrap(err)
			}
			return res.Idps, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return nil, nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
