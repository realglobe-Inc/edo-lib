package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoLoginRegistry(url, dbName, collName string, expiDur time.Duration) (LoginRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newLoginRegistry(base), nil
}
