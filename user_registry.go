package driver

import ()

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
