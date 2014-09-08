package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/url"
)

// 非キャッシュ用。
func NewWebServiceExplorer(addr string, ssl bool) (ServiceExplorer, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) ServiceUuid(servUri string) (servUuid string, err error) {
	resp, err := reg.Get(reg.prefix + "?service_uri=" + url.QueryEscape(servUri))
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var buff struct {
		Service struct {
			Uri string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
		return "", erro.Wrap(err)
	}
	return buff.Service.Uri, nil
}

// キャッシュ用。
func NewWebDatedServiceExplorer(addr string, ssl bool) (DatedServiceExplorer, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
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
					Uri string
				}
			}
			if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
				return "", nil, erro.Wrap(err)
			}
			return buff.Service.Uri, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return "", nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
