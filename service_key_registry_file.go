package driver

import (
	"time"
)

// バックエンドにファイルシステムを使う。

// 非キャッシュ用。
func NewFileServiceKeyRegistry(path string) ServiceKeyRegistry {
	return newServiceKeyRegistry(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}

// キャッシュ用。
func NewFileDatedServiceKeyRegistry(path string, expiDur time.Duration) DatedServiceKeyRegistry {
	return newDatedServiceKeyRegistry(newSynchronizedDatedKeyValueStore(newFileDatedKeyValueStore(path, expiDur)))
}
