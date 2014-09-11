package driver

import ()

func NewFileUserNameIndex(path string) UserNameIndex {
	return newUserNameIndex(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}
