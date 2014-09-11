package driver

import (
	"time"
)

// 期限が切れたら勝手に消える入れ物。
type TimeLimitedKeyValueStore interface {
	Get(key string) (value interface{}, err error)
	Put(key string, value interface{}, timLim time.Time) error
	Remove(key string) error
}
