package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoUserRegistry(url, dbName, collName string, expiDur time.Duration) (UserRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newUserRegistry(base), nil
}
