package driver

import (
	"io"
)

type RawDataStore interface {
	Get(key string, caStmp *Stamp) (data []byte, newCaStmp *Stamp, err error)
	Put(key string, data []byte) (*Stamp, error)
	Remove(key string) error

	io.Closer
}

// ListedRawDataStore に見せかけるための Keys を実装しない ListedRawDataStore もどき。
type listedRawDataStoreDummy struct {
	RawDataStore
}

func (drv *listedRawDataStoreDummy) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	panic("not implemented")
}
