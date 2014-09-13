package driver

import ()

// 非キャッシュ用。
func NewFileUserNameIndex(path string) UserNameIndex {
	return newUserNameIndex(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}
