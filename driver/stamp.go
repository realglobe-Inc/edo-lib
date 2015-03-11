package driver

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"net/http"
	"time"
)

// HTTP のキャッシュみたいなことができるように。
// 取得操作の場合、対象の更新日時がキャッシュされている更新日時以降、かつ、
// 対象のダイジェストがキャッシュされているダイジェストと異なる場合のみ、対象が返る。
//
//                          返り値 返りスタンプ
// 対象が無い               nil    nil
// キャッシュが古そうでない nil    あり
// キャッシュが古そう       あり   あり
//
// 更新操作の場合、新しいスタンプが返る。

// キャッシュの情報。
type Stamp struct {
	Date      time.Time `json:"date"                      bson:"date"`                      // 元データの更新日時。
	StaleDate time.Time `json:"stale_date,omitempty"      bson:"stale_date,omitempty"`      // 古くなったと疑わなければならない日時。
	ExpiDate  time.Time `json:"expiration_date,omitempty" bson:"expiration_date,omitempty"` // 廃棄しなければならない日時。
	Digest    string    `json:"digest"                    bson:"digest"`                    // ハッシュ値とか。
}

// キャッシュのタイムスタンプが対象のタイムスタンプより古そうなときのみ true。
func (caStmp *Stamp) Older(stmp *Stamp) bool {
	// 秒単位で比較するのは HTTP やら DB やらによって
	// 秒未満が切り捨てられることがあるから。
	if caDate, date := caStmp.Date.Unix(), stmp.Date.Unix(); caDate < date {
		// 日時が古い。
		return true
	} else if caDate == date && caStmp.Digest != stmp.Digest {
		// 日時は同じだがパッと見は違う。
		return true
	} else {
		// 日時が新しい、または、日時が同じでパッと見も同じ。
		return false
	}
}

func WriteStampToRequestHeader(stmp *Stamp, h http.Header) {
	h.Set("If-Modified-Since", stmp.Date.Format(http.TimeFormat))
	h.Set("If-None-Match", stmp.Digest)
}
func ParseStampFromRequestHeader(h http.Header) (*Stamp, error) {
	dateStr, dig := h.Get("If-Modified-Since"), h.Get("If-None-Match")
	if dig == "" && dateStr == "" {
		return nil, nil
	}

	stmp := &Stamp{Digest: dig}
	if dateStr != "" {
		var err error
		stmp.Date, err = time.Parse(http.TimeFormat, dateStr)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	return stmp, nil
}
func WriteStampToResponseHeader(stmp *Stamp, h http.Header) {
	h.Set("Last-Modified", stmp.Date.Format(http.TimeFormat))
	h.Set("Expires", stmp.ExpiDate.Format(http.TimeFormat))
	h.Set("ETag", stmp.Digest)
}
func ParseStampFromResponseHeader(h http.Header) (*Stamp, error) {
	dateStr, expiDateStr, dig := h.Get("Last-Modified"), h.Get("Expires"), h.Get("ETag")
	if dig == "" && dateStr == "" {
		return nil, nil
	}

	stmp := &Stamp{Digest: dig}
	if dateStr != "" {
		var err error
		stmp.Date, err = time.Parse(http.TimeFormat, dateStr)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	if expiDateStr != "" {
		var err error
		stmp.ExpiDate, err = time.Parse(http.TimeFormat, expiDateStr)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	// 有効期間の半分は無確認で用いる。
	now := time.Now()
	stmp.StaleDate = now.Add(stmp.ExpiDate.Sub(now) / 2)

	return stmp, nil
}
