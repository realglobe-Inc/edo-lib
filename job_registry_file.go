package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// data を JSON として、JobResult にデコードする。
func jobResultUnmarshal(data []byte) (value interface{}, err error) {
	var res JobResult
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &res, nil
}

// スレッドセーフ。
func NewFileJobRegistry(path string, expiDur time.Duration) JobRegistry {
	return newJobRegistry(NewFileTimeLimitedKeyValueStore(path, jsonKeyGen, json.Marshal, jobResultUnmarshal, expiDur))
}
