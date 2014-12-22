package driver

import ()

type KeyValueStore interface {
	Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error)
	Put(key string, val interface{}) (*Stamp, error)
	Remove(key string) error
}
