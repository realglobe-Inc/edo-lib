package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "user_attribute_key": "user-attribute-no-key",
//   "user_attribute": XXX
// }

// 非キャッシュ用。
func NewMongoUserAttributeRegistry(url, dbName, collName string) (UserAttributeRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, "user_attribute_key", "user_attribute")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newUserAttributeRegistry(base), nil
}
