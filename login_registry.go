package driver

import ()

// ユーザーの管理。
type LoginRegistry interface {
	User(accToken string) (usrUuid string, err error)
}
