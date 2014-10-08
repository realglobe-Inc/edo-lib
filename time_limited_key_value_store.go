package driver

import (
	"time"
)

type TimeLimitedKeyValueStore interface {
	Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error)
	Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error)
	Remove(key string) error
}
