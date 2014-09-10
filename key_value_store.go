package driver

import ()

type keyValueStore interface {
	get(key string) (value interface{}, err error)
	put(key string, value interface{}) error
	remove(key string) error
}

// キャッシュ用。
type datedKeyValueStore interface {
	stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error)
	stampedPut(key string, value interface{}) (*Stamp, error)
	remove(key string) error
}
