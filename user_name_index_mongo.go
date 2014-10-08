package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// スレッドセーフ。
func NewMongoUserNameIndex(url, dbName, collName string, expiDur time.Duration) (UserNameIndex, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newUserNameIndex(base), nil
}
