package driver

import ()

type RawDataStore interface {
	Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error)
	Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error)
	Put(key string, data []byte) (*Stamp, error)
	Remove(key string) error
}
