package driver

import ()

func NewWebUserNameIndex(prefix string) UserNameIndex {
	return newUserNameIndex(newWebKeyValueStore(prefix))
}
