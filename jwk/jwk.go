// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwk

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"github.com/realglobe-Inc/edo-lib/base64url"
	"github.com/realglobe-Inc/go-lib/erro"
	"math/big"
)

type Key interface {
	// kty
	Type() string
	// use
	Use() string
	// key_ops
	Operations() map[string]bool
	// alg
	Algorithm() string
	// kid
	Id() string

	Public() crypto.PublicKey
	Private() crypto.PrivateKey
	Common() []byte

	// json.Marshal すると JWK になるようなマップを返す。
	ToMap() map[string]interface{}
}

type keyImpl struct {
	kty string
	use string
	ops map[string]bool
	alg string
	id  string

	pub crypto.PublicKey
	pri crypto.PrivateKey
	com []byte
}

func New(rawKey interface{}, m map[string]interface{}) Key {
	key := baseFromMap(m)
	switch k := rawKey.(type) {
	case *ecdsa.PrivateKey:
		key.kty = ktyEc
		key.pri = k
		key.pub = &k.PublicKey
	case *ecdsa.PublicKey:
		key.kty = ktyEc
		key.pub = k
	case *rsa.PrivateKey:
		key.kty = ktyRsa
		key.pri = k
		key.pub = &k.PublicKey
	case *rsa.PublicKey:
		key.kty = ktyRsa
		key.pub = k
	case []byte:
		key.kty = ktyOct
		key.com = k
	default:
	}
	return key
}

// JWK を json.Unmarshal した結果からつくる。
func FromMap(m map[string]interface{}) (k Key, err error) {
	key := baseFromMap(m)
	switch key.kty {
	case ktyEc:
		key.pri, key.pub, err = ecFromMap(m)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	case ktyRsa:
		key.pri, key.pub, err = rsaFromMap(m)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	case ktyOct:
		key.com, err = commonFromMap(m)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	default:
		return nil, erro.New("kty " + key.kty + " is unsupported")
	}
	return key, nil
}

func baseFromMap(m map[string]interface{}) *keyImpl {
	var ops map[string]bool
	if opArray, _ := m[tagKey_ops].([]interface{}); len(opArray) > 0 {
		for _, opRaw := range opArray {
			if op, _ := opRaw.(string); op != "" {
				if ops == nil {
					ops = map[string]bool{}
				}
				ops[op] = true
			}
		}
	}

	key := &keyImpl{ops: ops}
	key.kty, _ = m[tagKty].(string)
	key.use, _ = m[tagUse].(string)
	key.alg, _ = m[tagAlg].(string)
	key.id, _ = m[tagKid].(string)
	return key
}

func ecFromMap(m map[string]interface{}) (crypto.PrivateKey, crypto.PublicKey, error) {
	pub, err := _ecPublicFromMap(m)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	var pri ecdsa.PrivateKey
	if dStr, _ := m["d"].(string); dStr == "" {
		// 公開鍵だった。
		return nil, pub, nil
	} else if dRaw, err := base64url.DecodeString(dStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else {
		pri.PublicKey = *pub
		pri.D = (&big.Int{}).SetBytes(dRaw)
	}

	return &pri, &pri.PublicKey, nil
}

func ecPublicFromMap(m map[string]interface{}) (crypto.PublicKey, error) {
	return _ecPublicFromMap(m)
}

func _ecPublicFromMap(m map[string]interface{}) (*ecdsa.PublicKey, error) {
	var pub ecdsa.PublicKey
	switch crv, _ := m[tagCrv].(string); crv {
	case "":
		return nil, erro.New("no crv")
	case "P-256":
		pub.Curve = elliptic.P256()
	case "P-384":
		pub.Curve = elliptic.P384()
	case "P-521":
		pub.Curve = elliptic.P521()
	default:
		return nil, erro.New("unsupoorted crv " + crv)
	}

	if xStr, _ := m["x"].(string); xStr == "" {
		return nil, erro.New("no x")
	} else if xRaw, err := base64url.DecodeString(xStr); err != nil {
		return nil, erro.Wrap(err)
	} else if yStr, _ := m["y"].(string); yStr == "" {
		return nil, erro.New("no y")
	} else if yRaw, err := base64url.DecodeString(yStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		pub.X = (&big.Int{}).SetBytes(xRaw)
		pub.Y = (&big.Int{}).SetBytes(yRaw)
	}

	return &pub, nil
}

func rsaFromMap(m map[string]interface{}) (crypto.PrivateKey, crypto.PublicKey, error) {
	pub, err := _rsaPublicFromMap(m)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	var pri rsa.PrivateKey
	if dStr, _ := m["d"].(string); dStr == "" {
		// 公開鍵だった。
		return nil, pub, nil
	} else if dRaw, err := base64url.DecodeString(dStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if pStr, _ := m["p"].(string); pStr == "" {
		return nil, nil, erro.New("no p")
	} else if pRaw, err := base64url.DecodeString(pStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if qStr, _ := m["q"].(string); qStr == "" {
		return nil, nil, erro.New("no q")
	} else if qRaw, err := base64url.DecodeString(qStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if dpStr, _ := m["dp"].(string); dpStr == "" {
		return nil, nil, erro.New("no dp")
	} else if dpRaw, err := base64url.DecodeString(dpStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if dqStr, _ := m["dq"].(string); dqStr == "" {
		return nil, nil, erro.New("no dq")
	} else if dqRaw, err := base64url.DecodeString(dqStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if qiStr, _ := m["qi"].(string); qiStr == "" {
		return nil, nil, erro.New("no qi")
	} else if qiRaw, err := base64url.DecodeString(qiStr); err != nil {
		return nil, nil, erro.Wrap(err)
	} else {
		pri.PublicKey = *pub
		pri.D = (&big.Int{}).SetBytes(dRaw)
		pri.Primes = []*big.Int{(&big.Int{}).SetBytes(pRaw), (&big.Int{}).SetBytes(qRaw)}
		pri.Precomputed.Dp = (&big.Int{}).SetBytes(dpRaw)
		pri.Precomputed.Dq = (&big.Int{}).SetBytes(dqRaw)
		pri.Precomputed.Qinv = (&big.Int{}).SetBytes(qiRaw)
		pri.Precomputed.CRTValues = []rsa.CRTValue{}
	}

	if crts, _ := m[tagOth].([]map[string]interface{}); crts != nil {
		for _, crt := range crts {
			if rStr, _ := crt["r"].(string); rStr == "" {
				return nil, nil, erro.New("no r")
			} else if rRaw, err := base64url.DecodeString(rStr); err != nil {
				return nil, nil, erro.Wrap(err)
			} else if dStr, _ := m["d"].(string); dStr == "" {
				return nil, nil, erro.New("no d")
			} else if dRaw, err := base64url.DecodeString(dStr); err != nil {
				return nil, nil, erro.Wrap(err)
			} else if tStr, _ := m["t"].(string); tStr == "" {
				return nil, nil, erro.New("no t")
			} else if tRaw, err := base64url.DecodeString(tStr); err != nil {
				return nil, nil, erro.Wrap(err)
			} else {
				pri.Precomputed.CRTValues = append(pri.Precomputed.CRTValues, rsa.CRTValue{
					R:     (&big.Int{}).SetBytes(rRaw),
					Exp:   (&big.Int{}).SetBytes(dRaw),
					Coeff: (&big.Int{}).SetBytes(tRaw),
				})
			}
		}
	}

	return &pri, pub, nil
}

func rsaPublicFromMap(m map[string]interface{}) (crypto.PublicKey, error) {
	return _rsaPublicFromMap(m)
}

func _rsaPublicFromMap(m map[string]interface{}) (*rsa.PublicKey, error) {
	var pub rsa.PublicKey
	if nStr, _ := m["n"].(string); nStr == "" {
		return nil, erro.New("no n")
	} else if nRaw, err := base64url.DecodeString(nStr); err != nil {
		return nil, erro.Wrap(err)
	} else if eStr, _ := m["e"].(string); eStr == "" {
		return nil, erro.New("no e")
	} else if eRaw, err := base64url.DecodeString(eStr); err != nil {
		return nil, erro.Wrap(err)
	} else {
		pub.N = (&big.Int{}).SetBytes(nRaw)
		pub.E = int((&big.Int{}).SetBytes(eRaw).Int64())
	}

	return &pub, nil
}

func commonFromMap(m map[string]interface{}) ([]byte, error) {
	if kStr, _ := m["k"].(string); kStr == "" {
		return nil, erro.New("no k")
	} else {
		return base64url.DecodeString(kStr)
	}
}

func (this *keyImpl) Type() string                { return this.kty }
func (this *keyImpl) Use() string                 { return this.use }
func (this *keyImpl) Operations() map[string]bool { return this.ops }
func (this *keyImpl) Algorithm() string           { return this.alg }
func (this *keyImpl) Id() string                  { return this.id }
func (this *keyImpl) Public() crypto.PublicKey    { return this.pub }
func (this *keyImpl) Private() crypto.PrivateKey  { return this.pri }
func (this *keyImpl) Common() []byte              { return this.com }

func (this *keyImpl) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		tagKty: this.kty,
	}

	if this.use != "" {
		m[tagUse] = this.use
	}
	if len(this.ops) > 0 {
		ops := []interface{}{}
		for op := range this.ops {
			ops = append(ops, op)
		}
		m[tagKey_ops] = ops
	}
	if this.alg != "" {
		m[tagAlg] = this.alg
	}
	if this.id != "" {
		m[tagKid] = this.id
	}

	switch this.kty {
	case ktyEc:
		if this.pri != nil {
			return ecToMap(this.pri.(*ecdsa.PrivateKey), m)
		} else {
			return ecPublicToMap(this.pub.(*ecdsa.PublicKey), m)
		}
	case ktyRsa:
		if this.pri != nil {
			return rsaToMap(this.pri.(*rsa.PrivateKey), m)
		} else {
			return rsaPublicToMap(this.pub.(*rsa.PublicKey), m)
		}
	case ktyOct:
		return commonToMap(this.com, m)
	default:
		return nil
	}
}

func rsaToMap(key *rsa.PrivateKey, m map[string]interface{}) map[string]interface{} {
	m = rsaPublicToMap(&key.PublicKey, m)
	m["d"] = base64url.EncodeToString(key.D.Bytes())
	m["p"] = base64url.EncodeToString(key.Primes[0].Bytes())
	m["q"] = base64url.EncodeToString(key.Primes[1].Bytes())
	m["dp"] = base64url.EncodeToString(key.Precomputed.Dp.Bytes())
	m["dq"] = base64url.EncodeToString(key.Precomputed.Dq.Bytes())
	m["qi"] = base64url.EncodeToString(key.Precomputed.Qinv.Bytes())
	if len(key.Precomputed.CRTValues) > 0 {
		var crts []map[string]interface{}
		for _, crt := range key.Precomputed.CRTValues {
			crts = append(crts, map[string]interface{}{
				"r": base64url.EncodeToString(crt.R.Bytes()),
				"d": base64url.EncodeToString(crt.Exp.Bytes()),
				"t": base64url.EncodeToString(crt.Coeff.Bytes()),
			})
		}
		m[tagOth] = crts
	}
	return m
}

func rsaPublicToMap(key *rsa.PublicKey, m map[string]interface{}) map[string]interface{} {
	m["n"] = base64url.EncodeToString(key.N.Bytes())
	m["e"] = base64url.EncodeToString(big.NewInt(int64(key.E)).Bytes())
	return m
}

func ecToMap(key *ecdsa.PrivateKey, m map[string]interface{}) map[string]interface{} {
	m = ecPublicToMap(&key.PublicKey, m)
	size := (key.Params().BitSize + 7) / 8
	m["d"] = base64url.EncodeToString(pad0(key.D.Bytes(), size))
	return m
}

func ecPublicToMap(key *ecdsa.PublicKey, m map[string]interface{}) map[string]interface{} {
	switch key.Params().BitSize {
	case 256:
		m[tagCrv] = "P-256"
	case 384:
		m[tagCrv] = "P-384"
	case 521:
		m[tagCrv] = "P-521"
	default:
		return nil
	}

	size := (key.Params().BitSize + 7) / 8
	m["x"] = base64url.EncodeToString(pad0(key.X.Bytes(), size))
	m["y"] = base64url.EncodeToString(pad0(key.Y.Bytes(), size))

	return m
}

func commonToMap(key []byte, m map[string]interface{}) map[string]interface{} {
	m["k"] = base64url.EncodeToString(key)
	return m
}

// size バイトになるように前に 0 を詰める。
func pad0(b []byte, size int) []byte {
	if len(b) < size {
		buff := make([]byte, size)
		copy(buff[size-len(b):], b)
		b = buff
	}
	return b
}
