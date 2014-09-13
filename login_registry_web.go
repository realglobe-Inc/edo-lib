package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebLoginRegistry(prefix string) LoginRegistry {
	return newWebDriver(prefix)
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
