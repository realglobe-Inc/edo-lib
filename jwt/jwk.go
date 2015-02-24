package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"github.com/realglobe-Inc/go-lib/erro"
	"math/big"
)

func KeyFromJwkMap(m map[string]interface{}) (key interface{}, err error) {
	switch kty, _ := m["kty"].(string); kty {
	case "oct":
		return commonKeyFromJwkMap(m)
	case "RSA":
		return rsaKeyFromJwkMap(m)
	case "EC":
		return ecdsaKeyFromJwkMap(m)
	default:
		return nil, erro.New("kty " + kty + " is unsupported")
	}
}

func commonKeyFromJwkMap(m map[string]interface{}) ([]byte, error) {
	if kStr, _ := m["k"].(string); kStr == "" {
		return nil, erro.New("no k")
	} else {
		return base64UrlDecodeString(kStr)
	}
}

func rsaKeyFromJwkMap(m map[string]interface{}) (interface{}, error) {
	pubKey, err := rsaPublicKeyFromJwkMap(m)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var priKey rsa.PrivateKey
	if dStr, _ := m["d"].(string); dStr == "" {
		// 公開鍵だった。
		return pubKey, nil
	} else if dRaw, err := base64UrlDecodeString(dStr); err != nil {
		return nil, erro.Wrap(err)
	} else if pStr, _ := m["p"].(string); pStr == "" {
		return nil, erro.New("no p")
	} else if pRaw, err := base64UrlDecodeString(pStr); err != nil {
		return nil, erro.Wrap(err)
	} else if qStr, _ := m["q"].(string); qStr == "" {
		return nil, erro.New("no q")
	} else if qRaw, err := base64UrlDecodeString(qStr); err != nil {
		return nil, erro.Wrap(err)
	} else if dpStr, _ := m["dp"].(string); dpStr == "" {
		return nil, erro.New("no dp")
	} else if dpRaw, err := base64UrlDecodeString(dpStr); err != nil {
		return nil, erro.Wrap(err)
	} else if dqStr, _ := m["dq"].(string); dqStr == "" {
		return nil, erro.New("no dq")
	} else if dqRaw, err := base64UrlDecodeString(dqStr); err != nil {
		return nil, erro.Wrap(err)
	} else if qiStr, _ := m["qi"].(string); qiStr == "" {
		return nil, erro.New("no qi")
	} else if qiRaw, err := base64UrlDecodeString(qiStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		priKey.PublicKey = *pubKey
		priKey.D = (&big.Int{}).SetBytes(dRaw)
		priKey.Primes = []*big.Int{(&big.Int{}).SetBytes(pRaw), (&big.Int{}).SetBytes(qRaw)}
		priKey.Precomputed.Dp = (&big.Int{}).SetBytes(dpRaw)
		priKey.Precomputed.Dq = (&big.Int{}).SetBytes(dqRaw)
		priKey.Precomputed.Qinv = (&big.Int{}).SetBytes(qiRaw)
		priKey.Precomputed.CRTValues = []rsa.CRTValue{}
	}

	if crts, _ := m["oth"].([]map[string]interface{}); crts != nil {
		for _, crt := range crts {
			if rStr, _ := crt["r"].(string); rStr == "" {
				return nil, erro.New("no r")
			} else if rRaw, err := base64UrlDecodeString(rStr); err != nil {
				return nil, erro.Wrap(err)
			} else if dStr, _ := m["d"].(string); dStr == "" {
				return nil, erro.New("no d")
			} else if dRaw, err := base64UrlDecodeString(dStr); err != nil {
				return nil, erro.Wrap(err)
			} else if tStr, _ := m["t"].(string); tStr == "" {
				return nil, erro.New("no t")
			} else if tRaw, err := base64UrlDecodeString(tStr); err != nil {
				return nil, erro.Wrap(err)
			} else {
				priKey.Precomputed.CRTValues = append(priKey.Precomputed.CRTValues, rsa.CRTValue{
					R:     (&big.Int{}).SetBytes(rRaw),
					Exp:   (&big.Int{}).SetBytes(dRaw),
					Coeff: (&big.Int{}).SetBytes(tRaw),
				})
			}
		}
	}

	return &priKey, nil
}

func rsaPublicKeyFromJwkMap(m map[string]interface{}) (*rsa.PublicKey, error) {
	var key rsa.PublicKey
	if nStr, _ := m["n"].(string); nStr == "" {
		return nil, erro.New("no n")
	} else if nRaw, err := base64UrlDecodeString(nStr); err != nil {
		return nil, erro.Wrap(err)
	} else if eStr, _ := m["e"].(string); eStr == "" {
		return nil, erro.New("no e")
	} else if eRaw, err := base64UrlDecodeString(eStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.N = (&big.Int{}).SetBytes(nRaw)
		key.E = int((&big.Int{}).SetBytes(eRaw).Int64())
	}
	return &key, nil
}

func ecdsaKeyFromJwkMap(m map[string]interface{}) (interface{}, error) {
	pubKey, err := ecdsaPublicKeyFromJwkMap(m)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var key ecdsa.PrivateKey
	if dStr, _ := m["d"].(string); dStr == "" {
		// 公開鍵だった。
		return pubKey, nil
	} else if dRaw, err := base64UrlDecodeString(dStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.PublicKey = *pubKey
		key.D = (&big.Int{}).SetBytes(dRaw)
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
		return nil, erro.New("crv " + crv + " is unsupoorted")
	}

	if xStr, _ := m["x"].(string); xStr == "" {
		return nil, erro.New("no x")
	} else if xRaw, err := base64UrlDecodeString(xStr); err != nil {
		return nil, erro.Wrap(err)
	} else if yStr, _ := m["y"].(string); yStr == "" {
		return nil, erro.New("no y")
	} else if yRaw, err := base64UrlDecodeString(yStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		key.X = (&big.Int{}).SetBytes(xRaw)
		key.Y = (&big.Int{}).SetBytes(yRaw)
	}
	return &key, nil
}

func KeyToJwkMap(key interface{}, m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = map[string]interface{}{}
	}

	switch k := key.(type) {
	case []byte:
		return commonKeyToJwkMap(k, m)
	case *rsa.PublicKey:
		return rsaPublicKeyToJwkMap(k, m)
	case *rsa.PrivateKey:
		return rsaPrivateKeyToJwkMap(k, m)
	case *ecdsa.PublicKey:
		return ecdsaPublicKeyToJwkMap(k, m)
	case *ecdsa.PrivateKey:
		return ecdsaPrivateKeyToJwkMap(k, m)
	default:
		return nil
	}
}

func commonKeyToJwkMap(key []byte, m map[string]interface{}) map[string]interface{} {
	m["kty"] = "oct"
	m["k"] = base64UrlEncodeToString(key)
	return m
}

func rsaPublicKeyToJwkMap(key *rsa.PublicKey, m map[string]interface{}) map[string]interface{} {
	m["kty"] = "RSA"
	m["n"] = base64UrlEncodeToString(key.N.Bytes())
	m["e"] = base64UrlEncodeToString(big.NewInt(int64(key.E)).Bytes())
	return m
}

func rsaPrivateKeyToJwkMap(key *rsa.PrivateKey, m map[string]interface{}) map[string]interface{} {
	m = rsaPublicKeyToJwkMap(&key.PublicKey, m)
	m["d"] = base64UrlEncodeToString(key.D.Bytes())
	m["p"] = base64UrlEncodeToString(key.Primes[0].Bytes())
	m["q"] = base64UrlEncodeToString(key.Primes[1].Bytes())
	m["dp"] = base64UrlEncodeToString(key.Precomputed.Dp.Bytes())
	m["dq"] = base64UrlEncodeToString(key.Precomputed.Dq.Bytes())
	m["qi"] = base64UrlEncodeToString(key.Precomputed.Qinv.Bytes())
	if len(key.Precomputed.CRTValues) > 0 {
		var crts []map[string]interface{}
		for _, crt := range key.Precomputed.CRTValues {
			crts = append(crts, map[string]interface{}{
				"r": base64UrlEncodeToString(crt.R.Bytes()),
				"d": base64UrlEncodeToString(crt.Exp.Bytes()),
				"t": base64UrlEncodeToString(crt.Coeff.Bytes()),
			})
		}
		m["oth"] = crts
	}
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
	m["x"] = base64UrlEncodeToString(paddedBigEndian(key.X, size))
	m["y"] = base64UrlEncodeToString(paddedBigEndian(key.Y, size))

	return m
}

func ecdsaPrivateKeyToJwkMap(key *ecdsa.PrivateKey, m map[string]interface{}) map[string]interface{} {
	m = ecdsaPublicKeyToJwkMap(&key.PublicKey, m)
	size := (key.Params().BitSize + 7) / 8
	m["d"] = base64UrlEncodeToString(paddedBigEndian(key.D, size))
	return m
}

func paddedBigEndian(val *big.Int, size int) []byte {
	b := val.Bytes()
	if len(b) < size {
		buff := make([]byte, size)
		copy(buff[size-len(b):], b)
		b = buff
	}
	return b
}
