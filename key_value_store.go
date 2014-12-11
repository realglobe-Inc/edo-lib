package driver

import ()

type KeyValueStore interface {
	Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error)
	Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error)
	Put(key string, value interface{}) (*Stamp, error)
	Remove(key string) error
}
