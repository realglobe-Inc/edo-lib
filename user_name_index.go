package driver

import ()

// ユーザー名からユーザー UUID を引く。
type UserNameIndex interface {
	UserUuid(usrName string) (usrUuid string, err error)
}

// 非キャッシュ用。
type userNameIndex struct {
	keyValueStore
}

func newUserNameIndex(base keyValueStore) *userNameIndex {
	return &userNameIndex{base}
}

func (reg *userNameIndex) UserUuid(usrName string) (usrUuid string, err error) {
	val, err := reg.get(usrName)
	if val != nil && val != "" {
		usrUuid = val.(string)
	}
	return usrUuid, err
}
