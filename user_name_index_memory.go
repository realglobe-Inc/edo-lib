package driver

import ()

// 非キャッシュ用。
type MemoryUserNameIndex struct {
	keyValueStore
}

func NewMemoryUserNameIndex() *MemoryUserNameIndex {
	return &MemoryUserNameIndex{newSynchronizedKeyValueStore(newMemoryKeyValueStore())}
}

func (reg *MemoryUserNameIndex) UserUuid(usrName string) (usrUuid string, err error) {
	val, err := reg.get(usrName)
	if val != nil && val != "" {
		usrUuid = val.(string)
	}
	return usrUuid, err
}

func (reg *MemoryUserNameIndex) AddUserUuid(usrName, usrUuid string) {
	reg.put(usrName, usrUuid)
}

func (reg *MemoryUserNameIndex) RemoveIdProvider(usrName string) {
	reg.remove(usrName)
}
