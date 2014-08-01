package driver

import (
	"time"
)

// ユーザーの管理。
type LoginRegistry interface {
	User(accToken string) (usrUuid string, err error)
}

// JavaScript 管理。
type JsRegistry interface {
	// オブジェクトのソースを取得する。
	Object(dir, objName string) (*Object, error)

	// オブジェクトのソースを登録する。一時的な利用を想定。
	AddObject(dir, objName string, obj *Object) error
	// オブジェクトのソースを削除する。
	RemoveObject(dir, objName string) error
}

type Object struct {
	Service bool `json:"service,omitempty"` // Web API からの利用を許可するか。
	Library bool `json:"library,omitempty"` // ライブラリとしての利用を許可するか。

	// 直接利用するライブラリとコネクタ。
	// Call するオブジェクトは含まない。
	// ライブラリは $library、コネクタは $$connector 形式で。
	Include []string `json:"include,omitempty"`

	Code string `json:"code"`
}

// ユーザー情報の管理。
type UserRegistry interface {
	// ユーザー情報を取得。
	Attributes(usrUuid string) (attrs map[string]interface{}, err error)
	Attribute(usrUuid, attrName string) (attr interface{}, err error)

	// ユーザー情報を変更。
	AddAttribute(usrUuid, attrName string, attr interface{}) error
	// ユーザー情報を削除。
	RemoveAttribute(usrUuid, attrName string) error
}

// ジョブ管理。
type JobRegistry interface {
	// 実行結果を取得する。
	Result(jobId string) (*JobResult, error)

	// 実行結果を登録する。
	AddResult(jobId string, res *JobResult, deadline time.Time) error
}

type JobResult struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

// 別名の管理。
// name_registry 参照。

// イベントの管理。
// event_registry 参照。

// サービス UUID の管理。
type ServiceRegistry interface {
	// EDO に登録されたサービスの管轄外向けエンドポイントから UUID を引く。
	Service(addr string) (servUuid string, err error)
}
