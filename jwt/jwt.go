package jwt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strings"
)

// JSON Web Token
type Jwt interface {
	Header(tag string) interface{}
	// val が nil や空文字列の場合は削除する。
	SetHeader(tag string, val interface{})
	Claim(tag string) interface{}
	// val が nil や空文字列の場合は削除する。
	SetClaim(tag string, val interface{})

	HeaderNames() map[string]bool
	ClaimNames() map[string]bool

	Encode() ([]byte, error)
	ToJson() (headJs, clmsJs []byte, err error)
}

func NewJwt() Jwt {
	return newJwt()
}

func ParseJwt(raw string) (Jwt, error) {
	jt := &jwt{}
	if err := jt.parse(raw); err != nil {
		return nil, erro.Wrap(err)
	}
	return jt, nil
}

type jwt struct {
	head map[string]interface{}
	clms map[string]interface{}

	// 以下キャッシュ。
	headJs []byte
	clmsJs []byte
}

func newJwt() *jwt {
	return &jwt{
		head: map[string]interface{}{},
		clms: map[string]interface{}{},
	}
}

func (this *jwt) parse(raw string) error {
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return erro.New("lack of JWT parts")
	} else if err := this.parseParts(parts[0], parts[1]); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func base64UrlDecodeString(s string) ([]byte, error) {
	if n := len(s) % 4; n == 2 {
		s += "=="
	} else if n == 3 {
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func base64UrlEncode(src []byte) []byte {
	buff := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(buff, src)
	return bytes.TrimRight(buff, "=")
}

func base64UrlEncodeToString(src []byte) string {
	return string(base64UrlEncode(src))
}

func (this *jwt) parseParts(headPart, clmsPart string) error {
	headJs, err := base64UrlDecodeString(headPart)
	if err != nil {
		return erro.Wrap(err)
	}
	var head map[string]interface{}
	if err := json.Unmarshal(headJs, &head); err != nil {
		return erro.Wrap(err)
	}

	clmsJs, err := base64UrlDecodeString(clmsPart)
	if err != nil {
		return erro.Wrap(err)
	}
	var clms map[string]interface{}
	if err := json.Unmarshal(clmsJs, &clms); err != nil {
		return erro.Wrap(err)
	}

	this.head = head
	this.clms = clms
	this.headJs = headJs
	this.clmsJs = clmsJs
	return nil
}

func (this *jwt) Header(tag string) interface{} {
	return this.head[tag]
}
func (this *jwt) SetHeader(tag string, val interface{}) {
	this.headJs = nil
	if val == nil || val == "" {
		delete(this.head, tag)
	} else {
		this.head[tag] = val
	}
}
func (this *jwt) Claim(clm string) interface{} {
	return this.clms[clm]
}
func (this *jwt) SetClaim(tag string, val interface{}) {
	this.clmsJs = nil
	if val == nil || val == "" {
		delete(this.clms, tag)
	} else {
		this.clms[tag] = val
	}
}
func (this *jwt) HeaderNames() map[string]bool {
	m := map[string]bool{}
	for k := range this.head {
		m[k] = true
	}
	return m
}
func (this *jwt) ClaimNames() map[string]bool {
	m := map[string]bool{}
	for k := range this.clms {
		m[k] = true
	}
	return m
}
func (this *jwt) Encode() ([]byte, error) {
	headJs, clmsJs, err := this.ToJson()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	buff := base64UrlEncode(headJs)
	buff = append(buff, '.')
	buff = append(buff, base64UrlEncode(clmsJs)...)
	return buff, nil
}

func (this *jwt) ToJson() (headJs, clmsJs []byte, err error) {
	if this.headJs == nil {
		this.headJs, err = json.Marshal(this.head)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
	}
	if this.clmsJs == nil {
		this.clmsJs, err = json.Marshal(this.clms)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
	}
	return this.headJs, this.clmsJs, nil
}

func Base64UrlDecodeString(s string) ([]byte, error) { return base64UrlDecodeString(s) }
func Base64UrlEncode(src []byte) []byte              { return base64UrlEncode(src) }
func Base64UrlEncodeToString(src []byte) string      { return base64UrlEncodeToString(src) }
