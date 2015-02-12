package driver

import ()

type ListedKeyValueStore interface {
	lister
	KeyValueStore
}
