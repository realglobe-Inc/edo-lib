package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// {
//     "user": {
//         "uuid": "user-no-uuid"
//     }
// }
func webUserUuidUnmarshal(data []byte) (interface{}, error) {
	var res struct {
		User struct {
			Uuid string
		}
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.User.Uuid, nil
}

type webUserNameIndex struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewWebUserNameIndex(prefix string) UserNameIndex {
	return &webUserNameIndex{NewWebKeyValueStore(prefix, nil, webUserUuidUnmarshal)}
}

func (reg *webUserNameIndex) UserUuid(usrName string, caStmp *Stamp) (usrUuid string, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrName, caStmp)
	if err != nil {
		return "", nil, erro.Wrap(err)
	} else if value == nil || value == "" {
		return "", newCaStmp, nil
	}
	return value.(string), newCaStmp, nil
}
