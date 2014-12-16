package driver

import ()

type KeyValueStore interface {
	Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error)
	Put(key string, value interface{}) (*Stamp, error)
	Remove(key string) error
}
