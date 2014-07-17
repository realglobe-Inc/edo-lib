package driver

import (
	"time"
)

// HTTP のキャッシュみたいなことができるように。
// 取得操作の場合、対象の最終更新日時がキャッシュのもの以降かつダイジェストがキャッシュのものと異なる場合のみ、
// 現在の対象とそのスタンプが返る。
// 更新操作の場合、対象のスタンプがキャッシュのものと等しい場合のみ操作が行われ、新しいスタンプが返る。

type stamp struct {
	lastDate time.Time
	expiTime time.Time
	digest   string
}

type jsRegistry interface {
	// オブジェクトのソースを取得する。
	object(dir, objName string, cachedStmp *stamp) (*Object, *stamp, error)
}

type userRegistry interface {
	// ユーザー情報を取得。
	attributes(usrUuid string, cachedStmp *stamp) (attrs map[string]interface{}, stmp *stamp, err error)
	attribute(usrUuid, attrName string, cachedStmp *stamp) (attr interface{}, stmp *stamp, err error)

	// ユーザー情報を変更。
	addAttribute(usrUuid, attrName string, attr interface{}, cachedStmp *stamp) (attrsStmp, attrStmp *stamp, err error)
	// ユーザー情報を削除。
	removeAttribute(usrUuid, attrName string, cachedStmp *stamp) (attrsStmp, attrStmp *stamp, err error)
}
