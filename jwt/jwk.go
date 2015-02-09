package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"math/big"
)

func PublicKeyFromJwkMap(m map[string]interface{}) (kid string, key crypto.PublicKey, err error) {
	kid, _ = m["kid"].(string)
	switch kty, _ := m["kty"].(string); kty {
	case "RSA":
		key, err = rsaPublicKeyFromJwkMap(m)
		if err != nil {
			return "", nil, erro.Wrap(err)
		}
	case "EC":
		key, err = ecdsaPublicKeyFromJwkMap(m)
		if err != nil {
			return "", nil, erro.Wrap(err)
		}
	default:
		return "", nil, nil
	}
	return kid, key, nil
}

func rsaPublicKeyFromJwkMap(m map[string]interface{}) (*rsa.PublicKey, error) {
	var key rsa.PublicKey

	if nStr, _ := m["n"].(string); nStr == "" {
		return nil, erro.New("no n")
	} else if nRaw, err := base64UrlDecodeString(nStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.N = (&big.Int{}).SetBytes(nRaw)
	}

	if eStr, _ := m["e"].(string); eStr == "" {
		return nil, erro.New("no e")
	} else if eRaw, err := base64UrlDecodeString(eStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.E = int((&big.Int{}).SetBytes(eRaw).Int64())
	}

	return &key, nil
}

func ecdsaPublicKeyFromJwkMap(m map[string]interface{}) (*ecdsa.PublicKey, error) {
	var key ecdsa.PublicKey

	switch crv, _ := m["crv"].(string); crv {
	case "":
		return nil, erro.New("no crv")
	case "P-256":
		key.Curve = elliptic.P256()
	case "P-384":
		key.Curve = elliptic.P384()
	case "P-521":
		key.Curve = elliptic.P521()
	default:
		return nil, erro.New("unsupported elliptic curve " + crv)
	}

	if xStr, _ := m["x"].(string); xStr == "" {
		return nil, erro.New("no x")
	} else if xRaw, err := base64UrlDecodeString(xStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.X = (&big.Int{}).SetBytes(xRaw)
	}
	if yStr, _ := m["y"].(string); yStr == "" {
		return nil, erro.New("no y")
	} else if yRaw, err := base64UrlDecodeString(yStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.Y = (&big.Int{}).SetBytes(yRaw)
	}
	return &key, nil
}

func PublicKeyToJwkMap(kid string, key crypto.PublicKey) map[string]interface{} {
	m := map[string]interface{}{}
	if kid != "" {
		m["kid"] = kid
	}
	switch k := key.(type) {
	case *rsa.PublicKey:
		return rsaPublicKeyToJwkMap(k, m)
	case *ecdsa.PublicKey:
		return ecdsaPublicKeyToJwkMap(k, m)
	default:
		return nil
	}
}

func rsaPublicKeyToJwkMap(key *rsa.PublicKey, m map[string]interface{}) map[string]interface{} {
	m["kty"] = "RSA"
	m["n"] = base64UrlEncodeToString(key.N.Bytes())
	m["e"] = base64UrlEncodeToString(big.NewInt(int64(key.E)).Bytes())
	return m
}

func ecdsaPublicKeyToJwkMap(key *ecdsa.PublicKey, m map[string]interface{}) map[string]interface{} {
	m["kty"] = "EC"

	switch key.Params().BitSize {
	case 256:
		m["crv"] = "P-256"
	case 384:
		m["crv"] = "P-384"
	case 521:
		m["crv"] = "P-521"
	default:
		return nil
	}

	size := (key.Params().BitSize + 7) / 8

	xBuff := key.X.Bytes()
	if len(xBuff) < size {
		xBuff = append(make([]byte, size-len(xBuff)), xBuff...)
	}
	m["x"] = base64UrlEncodeToString(xBuff)

	yBuff := key.Y.Bytes()
	if len(yBuff) < size {
		yBuff = append(make([]byte, size-len(yBuff)), yBuff...)
	}
	m["y"] = base64UrlEncodeToString(yBuff)

	return m
}
