package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// mongodb をバックエンドに使う。
// {
//   "id_provider_uuid": "aaaa-bbbb-cccc",
//   "id_provider_query_uri": "https://realglobe.jp/query"
// }

// 非キャッシュ用。
func NewMongoIdProviderRegistry(url, dbName, collName string) (IdProviderRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, "id_provider_uuid", "id_provider_query_uri")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newIdProviderRegistry(base), nil
}

// キャッシュ用。
func NewMongoDatedIdProviderRegistry(url, dbName, collName string, expiDur time.Duration) (DatedIdProviderRegistry, error) {
	base, err := newMongoDatedKeyValueStore(url, dbName, collName, expiDur, "id_provider_uuid", "id_provider_query_uri")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newDatedIdProviderRegistry(base), nil
}
