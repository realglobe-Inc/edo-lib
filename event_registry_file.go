package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

func jsonKeyGen(before string) string {
	return before + ".json"
}

// スレッドセーフ。
func NewFileEventRegistry(path string, expiDur time.Duration) EventRegistry {
	return newEventRegistry(NewFileKeyValueStore(path, jsonKeyGen, json.Marshal, func(data []byte) (interface{}, error) {
		var eventToHndl map[string]Handler
		if err := json.Unmarshal(data, &eventToHndl); err != nil {
			return nil, erro.Wrap(err)
		}
		return eventToHndl, nil
	}, expiDur))
}
