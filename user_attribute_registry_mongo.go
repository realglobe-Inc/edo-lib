package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoUserAttributeRegistry(url, dbName, collName string, expiDur time.Duration) (UserAttributeRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newUserAttributeRegistry(base), nil
}
