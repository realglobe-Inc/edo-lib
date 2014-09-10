package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// mongodb をバックエンドに使う。
// {
//   "service_uuid": "aaaa-bbbb-cccc",
//   "service_public_key":  "XXXXX"
// }

// 非キャッシュ用。
func NewMongoServiceKeyRegistry(url, dbName, collName string) (ServiceKeyRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, "service_uuid", "service_public_key")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newServiceKeyRegistry(base), nil
}

// キャッシュ用。
func NewMongoDatedServiceKeyRegistry(url, dbName, collName string, expiDur time.Duration) (DatedServiceKeyRegistry, error) {
	base, err := newMongoDatedKeyValueStore(url, dbName, collName, expiDur, "service_uuid", "service_public_key")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newDatedServiceKeyRegistry(base), nil
}
