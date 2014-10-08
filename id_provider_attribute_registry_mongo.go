package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoIdProviderAttributeRegistry(url, dbName, collName string, expiDur time.Duration) (IdProviderAttributeRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newIdProviderAttributeRegistry(base), nil
}
