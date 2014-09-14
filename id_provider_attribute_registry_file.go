package driver

import (
	"time"
)

// バックエンドにファイルシステムを使う。

// 非キャッシュ用。
func NewFileIdProviderAttributeRegistry(path string) IdProviderAttributeRegistry {
	return newIdProviderAttributeRegistry(newSynchronizedKeyValueStore(newFileKeyValueStore(path)))
}

// キャッシュ用。
func NewFileDatedIdProviderAttributeRegistry(path string, expiDur time.Duration) DatedIdProviderAttributeRegistry {
	return newDatedIdProviderAttributeRegistry(newSynchronizedDatedKeyValueStore(newFileDatedKeyValueStore(path, expiDur)))
}
