package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebIdProviderRegistry(addr string, ssl bool) (IdProviderRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	resp, err := reg.Get(reg.prefix)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	var buff struct {
		Id_provider struct {
			Query_uri string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
		return "", erro.Wrap(err)
	}
	return buff.Id_provider.Query_uri, nil
}

// キャッシュ用。
func NewWebDatedIdProviderRegistry(addr string, ssl bool) (DatedIdProviderRegistry, error) {
	return newWebDriver(addr, ssl)
}

func (reg *webDriver) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
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
				Id_provider struct {
					Query_uri string
				}
			}
			if err := json.NewDecoder(resp.Body).Decode(&buff); err != nil {
				return "", nil, erro.Wrap(err)
			}
			return buff.Id_provider.Query_uri, stmp, nil
		default:
			panic("implementation error.")
		}
	default:
		return "", nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
}
