package driver

import (
	"time"
)

// バックエンドにファイルシステムを使う。

// 非キャッシュ用。
func NewFileIdProviderRegistry(path string) IdProviderRegistry {
	return newIdProviderRegistry(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}

// キャッシュ用。
func NewFileDatedIdProviderRegistry(path string, expiDur time.Duration) DatedIdProviderRegistry {
	return newDatedIdProviderRegistry(newSynchronizedDatedKeyValueStore(newFileDatedKeyValueStore(path, expiDur)))
}
