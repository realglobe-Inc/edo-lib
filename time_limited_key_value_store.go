package driver

import (
	"time"
)

type TimeLimitedKeyValueStore interface {
	Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error)
	Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error)
	Remove(key string) error
}
