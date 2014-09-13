package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
)

// 非キャッシュ用。
func NewWebNameRegistry(prefix string) NameRegistry {
	return newWebDriver(prefix)
}

func (reg *webDriver) Address(name string) (addr string, err error) {
	resp, err := reg.Get(reg.prefix + "/node/" + name)
	if err != nil {
		return "", erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode != http.StatusOK {
		return "", erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addr); err != nil {
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *webDriver) Addresses(name string) (addrs []string, err error) {
	resp, err := reg.Get(reg.prefix + "/tree/" + name)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, erro.New("invalid status ", resp.StatusCode, " "+http.StatusText(resp.StatusCode)+".")
	}
	if err := json.NewDecoder(resp.Body).Decode(&addrs); err != nil {
		return nil, erro.Wrap(err)
	}
	return addrs, nil
}
