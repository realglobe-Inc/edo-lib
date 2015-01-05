package util

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
	SetHeader(tag string, val interface{})
	Claim(tag string) interface{}
	SetClaim(tag string, val interface{})

	Encode() ([]byte, error)
}

func NewJwt() Jwt {
	return newJwt()
}

func ParseJwt(raw string) (Jwt, error) {
	parts := strings.Split(raw, ".")
	if len(parts) < 2 {
		return nil, erro.New("lack of JWT parts")
	}
	return parseJwt(parts[0], parts[1])
}

type jwt struct {
	head map[string]interface{}
	clms map[string]interface{}
}

func newJwt() *jwt {
	return &jwt{
		head: map[string]interface{}{},
		clms: map[string]interface{}{},
	}
}

func base64UrlDecodeString(str string) ([]byte, error) {
	if n := len(str) % 4; n == 2 {
		str += "=="
	} else if n == 3 {
		str += "="
	}
	return base64.URLEncoding.DecodeString(str)
}

func base64UrlEncode(input []byte) []byte {
	buff := make([]byte, base64.URLEncoding.EncodedLen(len(input)))
	base64.URLEncoding.Encode(buff, input)
	return bytes.TrimRight(buff, "=")
}

func base64UrlEncodeToString(input []byte) string {
	return string(base64UrlEncode(input))
}

func parseJwt(headPart, clmsPart string) (*jwt, error) {
	headJson, err := base64UrlDecodeString(headPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	var head map[string]interface{}
	if err := json.Unmarshal(headJson, &head); err != nil {
		return nil, erro.Wrap(err)
	}

	clmsJson, err := base64UrlDecodeString(clmsPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	var clms map[string]interface{}
	if err := json.Unmarshal(clmsJson, &clms); err != nil {
		return nil, erro.Wrap(err)
	}

	return &jwt{head, clms}, nil
}

func (this *jwt) Header(tag string) interface{} {
	return this.head[tag]
}
func (this *jwt) SetHeader(tag string, val interface{}) {
	this.head[tag] = val
}
func (this *jwt) Claim(clm string) interface{} {
	return this.clms[clm]
}
func (this *jwt) SetClaim(tag string, val interface{}) {
	this.clms[tag] = val
}
func (this *jwt) Encode() ([]byte, error) {
	headJson, err := json.Marshal(this.head)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	clmsJson, err := json.Marshal(this.clms)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	buff := base64UrlEncode(headJson)
	buff = append(buff, '.')
	buff = append(buff, base64UrlEncode(clmsJson)...)

	return buff, nil
}
