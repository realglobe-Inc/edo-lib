package driver

import (
	"io"
)

type Lister interface {
	lister
	io.Closer
}

type lister interface {
	Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error)
}
