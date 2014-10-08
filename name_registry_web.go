package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// data を JSON として、[]string にデコードする。
func stringArrayUnmarshal(data []byte) (interface{}, error) {
	var res []string
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

type webNameRegistry struct {
	addr  KeyValueStore
	addrs KeyValueStore
}

// スレッドセーフ。
func NewWebNameRegistry(prefix string) NameRegistry {
	return &webNameRegistry{
		NewWebKeyValueStore(prefix, nil, jsonUnmarshal),
		NewWebKeyValueStore(prefix, nil, stringArrayUnmarshal),
	}
}

func (reg *webNameRegistry) Address(name string, caStmp *Stamp) (addr string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.addr.Get("node/"+name, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}

func (reg *webNameRegistry) Addresses(name string, caStmp *Stamp) (addrs []string, newCaStmp *Stamp, err error) {
	value, _, err := reg.addr.Get("tree/"+name, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.([]string), newCaStmp, nil
}
