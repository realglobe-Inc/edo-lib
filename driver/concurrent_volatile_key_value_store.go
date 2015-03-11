package driver

import (
	"time"
)

// 並列使用時に便利なメソッドを持つ制限時間付きデータ用コンテナ。
type ConcurrentVolatileKeyValueStore interface {
	VolatileKeyValueStore

	// 以下、key と eKey の範囲が被っていた場合の挙動は保証しない。
	// エントリを返す。
	Entry(eKey string) (eVal string, err error)
	// エントリを設定する。
	SetEntry(eKey, eVal string, eExpDate time.Time) error
	// エントリを設定しつつ、値を返す。
	GetAndSetEntry(key string, caStmp *Stamp, eKey, eVal string, eExpiDate time.Time) (val interface{}, newCaStmp *Stamp, err error)
	// eVal がエントリに設定されていれば、値を設定する。
	PutIfEntered(key string, val interface{}, expiDate time.Time, eKey, eVal string) (entered bool, newCaStmp *Stamp, err error)
}
