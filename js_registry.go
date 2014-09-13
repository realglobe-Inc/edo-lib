package driver

import ()

type Object struct {
	Service bool `json:"service,omitempty"` // Web API からの利用を許可するか。
	Library bool `json:"library,omitempty"` // ライブラリとしての利用を許可するか。

	// 直接利用するライブラリとコネクタ。
	// Call するオブジェクトは含まない。
	// ライブラリは $library、コネクタは $$connector 形式で。
	Include []string `json:"include,omitempty"`

	Code string `json:"code"`
}

// 非キャッシュ用。
type JsRegistry interface {
	// オブジェクトのソースを取得する。
	Object(dir, objName string) (*Object, error)

	// オブジェクトのソースを登録する。一時的な利用を想定。
	AddObject(dir, objName string, obj *Object) error
	// オブジェクトのソースを削除する。
	RemoveObject(dir, objName string) error
}

// キャッシュ用。
type JsBackend interface {
	// オブジェクトのソースを取得する。
	StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error)
}

type JsBackendRegistry interface {
	JsRegistry
	JsBackend
}
