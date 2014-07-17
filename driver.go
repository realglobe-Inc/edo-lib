package driver

import (
	"fmt"
	"time"
)

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

// ユーザーの管理。
type LoginRegistry interface {
	User(accToken string) (usrUuid string, err error)
}

// ユーザー情報の管理。
type UserRegistry interface {
	// ユーザー情報を取得。
	Attributes(usrUuid string) (map[string]interface{}, error)
	Attribute(usrUuid, attrName string) (interface{}, error)

	// ユーザー情報を変更。
	AddAttribute(usrUuid, attrName string, attr interface{}) error
	// ユーザー情報を削除。
	RemoveAttribute(usrUuid, attrName string) error
}

// ジョブ管理。
type JobRegistry interface {
	// 実行結果を取得する。
	Result(usrUuid string, jobId uint64) (res interface{}, err error)

	// 実行結果を登録する。
	AddResult(usrUuid string, jobId uint64, res interface{}, deadline time.Time) error
}

// 別名の管理。
type NameRegistry interface {
	// アドレスを引く。
	Address(name string) (addr string, err error)
	// name はドメイン形式（. 区切りで後ろが親）の木構造のノードを表し、そのノード以下の部分木に含まれる全てのアドレスを返す。
	Addresses(name string) (addrs []string, err error)
}

// イベントの管理。
type EventRegistry interface {
	// ハンドラを取得する。
	// イベントは / 区切りで木構造のノードを表し、そのノード以下の部分木に含まれる全てのハンドラを返す。
	Handler(usrUuid, event string) (Handler, error)

	// ハンドラを登録する。
	AddHandler(usrUuid, event string, hndl Handler) error
	// ハンドラを削除する。
	RemoveHandler(usrUuid, event string) error
}

type Handler []*HandlerElement

type HandlerElement struct {
	Url     string            `json:"url"`
	Method  string            `json:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`

	Rules []*HandlerRule `json:"rules,omitempty"`
}

func (elem *HandlerElement) String() string {
	return fmt.Sprint("{"+elem.Url+" "+elem.Method+" ", elem.Headers, " ", len(elem.Body), "}")
}

// イベントの付属パラメータの扱いを記述する予定。
type HandlerRule struct {
}

// イベントの処理。
type EventRouter interface {
	// イベントを発生させる。
	Fire(usrUuid, event string, body interface{}) error
}
