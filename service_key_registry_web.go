package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebServiceKeyRegistry(addr string, ssl bool) (ServiceKeyRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) ServiceKey(servUuid string) (key string, err error) {
	resp, err := reg.Get(reg.prefix)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var buff struct {
		Service struct {
			Public_key string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
		return "", erro.Wrap(err)
	}
	return buff.Service.Public_key, nil
}

// キャッシュ用。
func NewWebDatedServiceKeyRegistry(addr string, ssl bool) (DatedServiceKeyRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error) {
	req, err := http.NewRequest("GET", reg.prefix, nil)
	if caStmp != nil {
		WriteStampToRequestHeader(caStmp, req.Header)
	}

	util.LogRequest(req, true)
	resp, err := reg.Do(req)
	if err != nil {
		return "", nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	util.LogResponse(resp, true)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return "", nil, nil
	case http.StatusNotModified, http.StatusOK:
		stmp, err := ParseStampFromResponseHeader(resp.Header)
		if err != nil {
			return "", nil, erro.Wrap(err)
		}

		switch resp.StatusCode {
		case http.StatusNotModified:
			return "", stmp, nil
		case http.StatusOK:
			var buff struct {
				Service struct {
					Public_key string
				}
			}
			if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
				return "", nil, erro.Wrap(err)
			}
			return buff.Service.Public_key, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return "", nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
