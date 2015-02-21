package jwt

import (
	"encoding/base64"
)

// 末尾に = を足さない Base64URL エンコード。

func base64UrlDecodeString(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func base64UrlEncode(src []byte) []byte {
	buff := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(buff, src)
	switch len(src) % 3 {
	case 1:
		return buff[:len(buff)-2]
	case 2:
		return buff[:len(buff)-1]
	default:
		return buff
	}
}

func base64UrlEncodeToString(src []byte) string {
	return string(base64UrlEncode(src))
}

func Base64UrlDecodeString(s string) ([]byte, error) { return base64UrlDecodeString(s) }
func Base64UrlEncode(src []byte) []byte              { return base64UrlEncode(src) }
func Base64UrlEncodeToString(src []byte) string      { return base64UrlEncodeToString(src) }
