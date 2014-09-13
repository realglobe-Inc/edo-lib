package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//   "user": {
//     "uuid": "abcd-no-uuid"
//   }
// }

// 非キャッシュ用。
func NewWebUserNameIndex(prefix string) UserNameIndex {
	return newWebUserNameIndex(newWebKeyValueStore(prefix))
}

type webUserNameIndex struct {
	keyValueStore
}

func newWebUserNameIndex(base keyValueStore) *webUserNameIndex {
	return &webUserNameIndex{base}
}

func (reg *webUserNameIndex) UserUuid(usrName string) (usrUuid string, err error) {
	val, err := reg.get(usrName)
	if err != nil {
		return "", erro.Wrap(err)
	} else if val == nil || val == "" {
		return "", nil
	}
	return val.(string), nil
}
