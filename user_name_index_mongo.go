package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

func NewMongoUserNameIndex(url, dbName, collName string) (UserNameIndex, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, "user_name", "user_uuid")
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newUserNameIndex(base), nil
}
