package driver

import ()

type ListedKeyValueStore interface {
	Lister
	KeyValueStore
}
