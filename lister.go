package driver

import ()

type Lister interface {
	Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error)
}
