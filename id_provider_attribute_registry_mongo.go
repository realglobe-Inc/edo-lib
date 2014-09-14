package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// mongodb をバックエンドに使う。
// {
//   "id_provider_attribute_key": "id-provider-attribute-no-key",
//   "id_provider_attribute": XXX
// }

// 非キャッシュ用。
func NewMongoIdProviderAttributeRegistry(url, dbName, collName string) (IdProviderAttributeRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, "id_provider_attribute_key", "id_provider_attribute")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newIdProviderAttributeRegistry(base), nil
}

// キャッシュ用。
func NewMongoDatedIdProviderAttributeRegistry(url, dbName, collName string, expiDur time.Duration) (DatedIdProviderAttributeRegistry, error) {
	base, err := newMongoDatedKeyValueStore(url, dbName, collName, expiDur, "id_provider_attribute_key", "id_provider_attribute")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newDatedIdProviderAttributeRegistry(base), nil
}
