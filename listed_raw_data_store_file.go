package driver

import (
	"time"
)

type fileListedRawDataStore struct {
	Lister
	RawDataStore
}

// スレッドセーフ。
func NewFileListedRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) ListedRawDataStore {
	return newSynchronizedListedRawDataStore(newCachingListedRawDataStore(newFileListedRawDataStore(path, keyToPath, pathToKey, staleDur, expiDur)))
}

// スレッドセーフではない。
func newFileListedRawDataStore(path string, keyToPath, pathToKey func(string) string, staleDur, expiDur time.Duration) *fileListedRawDataStore {
	return &fileListedRawDataStore{
		newFileLister(path, pathToKey, staleDur, expiDur),
		newFileRawDataStore(path, keyToPath, staleDur, expiDur),
	}
}
