package driver

import ()

// 中央レジストリに接続する本格的なレジストリ。
// キャッシュの管理も行う。

func NewBkfJsRegistry(addr string, ssl bool, secretKey string) (JsRegistry, error) {
	panic("not implemented.")
}

func NewBkfUserRegistry(addr string, ssl bool, secretKey string) (UserRegistry, error) {
	panic("not implemented.")
}

func NewBkfNameRegistry(addr string, ssl bool, secretKey string) (NameRegistry, error) {
	panic("not implemented.")
}

func NewBkfEventRegistry(addr string, ssl bool, secretKey string) (EventRegistry, error) {
	panic("not implemented.")
}
