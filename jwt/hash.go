package jwt

import (
	"crypto/sha256"
	"crypto/sha512"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"hash"
)

// JWS の alg ヘッダに対応するハッシュ関数を返す。
func HashFunction(alg string) (hash.Hash, error) {
	switch alg {
	case "none":
		return nil, nil
	case "HS256", "RS256", "ES256", "PS256":
		return sha256.New(), nil
	case "HS384", "RS384", "ES384", "PS384":
		return sha512.New384(), nil
	case "HS512", "RS512", "ES512", "PS512":
		return sha512.New(), nil
	default:
		return nil, erro.New("alg " + alg + " is unsupported")
	}
}
