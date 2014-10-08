package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

func handlerUnmarshal(data []byte) (value interface{}, err error) {
	var res Handler
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res, nil
}

// サーバーからイベントハンドラのリストを受け取るだけ。
type webEventRegistry struct {
	base KeyValueStore
}

func NewWebEventRegistry(prefix string) EventRegistry {
	return &webEventRegistry{NewWebKeyValueStore(prefix, json.Marshal, handlerUnmarshal)}
}

func (reg *webEventRegistry) Handler(usrUuid, event string, caStmp *Stamp) (hndl Handler, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrUuid+event, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}
	return value.(Handler), newCaStmp, nil
}

func (reg *webEventRegistry) AddHandler(usrUuid, event string, hndl Handler) (newCaStmp *Stamp, err error) {
	return reg.base.Put(usrUuid+event, hndl)
}

func (reg *webEventRegistry) RemoveHandler(usrUuid, event string) error {
	return reg.base.Remove(usrUuid + event)
}
