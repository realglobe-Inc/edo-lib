package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"time"
)

// HTTP のキャッシュみたいなことができるように。
// 取得操作の場合、対象の更新日時がキャッシュの日時以降、または、ダイジェストがキャッシュと異なる場合のみ、現在の対象が返る。
//
// 対象が            返り値 返りスタンプ
// 無い              nil    nil
// キャッシュと同じ  nil    あり
// キャッシュと違う  あり   あり
//
// 更新操作の場合、対象がキャッシュと等しい場合のみ操作が行われ、新しいスタンプが返る。

// キャッシュの情報。
type Stamp struct {
	Date     time.Time `json:"date"                      bson:"date"`                      // 元データの更新日時。
	ExpiDate time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"` // 有効期限。
	Digest   string    `json:"digest"                    bson:"digest"`                    // ハッシュ値とか。
}

func ParseStampFromRequestHeader(h http.Header) (*Stamp, error) {
	dig, dateStr := h.Get("If-None-Match"), h.Get("If-Modified-Since")
	if dig == "" && dateStr == "" {
		return nil, nil
	}
	stmp := &Stamp{}
	stmp.Digest = dig
	if dateStr != "" {
		date, err := time.Parse(time.RFC1123, dateStr)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		stmp.Date = date
	}
	return stmp, nil
}
func WriteStampToResponseHeader(stmp *Stamp, h http.Header) {
	h.Set("Last-Modified", stmp.Date.Format(time.RFC1123))
	h.Set("Expires", stmp.ExpiDate.Format(time.RFC1123))
	h.Set("ETag", stmp.Digest)
}

type DatedIdProviderLister interface {
	StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error)
}

type JsBackend interface {
	// オブジェクトのソースを取得する。
	StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error)
}

type JsBackendRegistry interface {
	JsRegistry
	JsBackend
}
