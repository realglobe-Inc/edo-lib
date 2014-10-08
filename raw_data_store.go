package driver

import ()

type RawDataStore interface {
	Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error)
	Put(key string, data []byte) (*Stamp, error)
	Remove(key string) error
}
