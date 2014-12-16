package driver

import ()

type ListedRawDataStore interface {
	Lister
	RawDataStore
}
