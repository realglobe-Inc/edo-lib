package crypto

import (
	"crypto"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"strings"
)

var hashToStr = map[crypto.Hash]string{
	crypto.MD4:       "MD4",
	crypto.MD5:       "MD5",
	crypto.SHA1:      "SHA1",
	crypto.SHA224:    "SHA224",
	crypto.SHA256:    "SHA256",
	crypto.SHA384:    "SHA384",
	crypto.SHA512:    "SHA512",
	crypto.MD5SHA1:   "MD5SHA1",
	crypto.RIPEMD160: "RIPEMD160",
	crypto.SHA3_224:  "SHA3_224",
	crypto.SHA3_256:  "SHA3_256",
	crypto.SHA3_384:  "SHA3_384",
	crypto.SHA3_512:  "SHA3_512",
}

var strToHash = map[string]crypto.Hash{}

func init() {
	for h, s := range hashToStr {
		strToHash[s] = h
	}
}

// 入力の大文字、小文字は区別しない。
func ParseHashFunction(s string) (crypto.Hash, error) {
	h, ok := strToHash[strings.ToUpper(s)]
	if !ok {
		return 0, erro.New("hash " + s + " is unsupported")
	}
	return h, nil
}

func HashFunctionString(h crypto.Hash) string {
	s, ok := hashToStr[h]
	if !ok {
		return "unknown"
	}
	return s
}
