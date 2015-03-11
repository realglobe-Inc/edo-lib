package driver

import ()

type ListedRawDataStore interface {
	lister
	RawDataStore
}
