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

package jwt

import (
	"bytes"
	"crypto"
	_ "crypto/ecdsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/json"
	"github.com/realglobe-Inc/edo-lib/base64url"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"github.com/realglobe-Inc/edo-lib/secrand"
	"github.com/realglobe-Inc/go-lib/erro"
)

// JSON Web Token
type Jwt struct {
	// 生のヘッダ。`{"alg":"ES256"}` とか。
	rawHead []byte
	// ヘッダ。
	head map[string]interface{}

	// 生の本文。`{"iss":"https://example.org"}` とかを想定するが、そうでなくても良い。
	rawBody []byte
	// クレームセット。
	clms map[string]interface{}

	// JWS の諸々。
	sig          []byte
	headBodyPart []byte
	// JWE の諸々。
	encedKey, initVec, enced, authTag []byte
	headPart                          []byte
	// Compact serialization.
	compact []byte
}

func New() *Jwt {
	return &Jwt{}
}

// 中間生成物を削除する。
func (this *Jwt) clear() {
	this.sig = nil
	this.enced = nil
	this.compact = nil
}

func (this *Jwt) SetRawHeader(raw []byte) {
	this.rawHead = raw
	this.head = nil
	this.clear()
}

func (this *Jwt) SetRawBody(raw []byte) {
	this.rawBody = raw
	this.clms = nil
	this.clear()
}

// val が nil なら削除する。
func (this *Jwt) SetHeader(tag string, val interface{}) {
	if this.head == nil {
		this.head = parseJsonOrNew(this.rawHead)
	}
	if val == nil {
		delete(this.head, tag)
	} else {
		this.head[tag] = val
	}
	this.rawHead = nil
	this.clear()
}

// val が nil なら削除する。
func (this *Jwt) SetClaim(tag string, val interface{}) {
	if this.clms == nil {
		this.clms = parseJsonOrNew(this.rawBody)
	}
	if val == nil {
		delete(this.clms, tag)
	} else {
		this.clms[tag] = val
	}
	this.rawBody = nil
	this.clear()
}

func (this *Jwt) Sign(keys []jwk.Key) error {
	if this.sig != nil {
		return nil
	}

	rawHead, err := this.getRawHeader()
	if err != nil {
		return erro.Wrap(err)
	}
	rawBody, err := this.getRawBody()
	if err != nil {
		return erro.Wrap(err)
	}

	headBodyPart := base64url.Encode(rawHead)
	headBodyPart = append(headBodyPart, '.')
	headBodyPart = append(headBodyPart, base64url.Encode(rawBody)...)

	kid, _ := this.Header(tagKid).(string)
	alg, _ := this.Header(tagAlg).(string)

	var sig []byte
	switch alg {
	case algNone:
		sig = []byte{}
	case algHs256:
		sig, err = hsSign(findKey(keys, kid, ktyOct, useSig, opSign, alg), crypto.SHA256, headBodyPart)
	case algHs384:
		sig, err = hsSign(findKey(keys, kid, ktyOct, useSig, opSign, alg), crypto.SHA384, headBodyPart)
	case algHs512:
		sig, err = hsSign(findKey(keys, kid, ktyOct, useSig, opSign, alg), crypto.SHA512, headBodyPart)
	case algRs256:
		sig, err = rsSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA256, headBodyPart)
	case algRs384:
		sig, err = rsSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA384, headBodyPart)
	case algRs512:
		sig, err = rsSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA512, headBodyPart)
	case algEs256:
		sig, err = esSign(findKey(keys, kid, ktyEc, useSig, opSign, alg), crypto.SHA256, headBodyPart)
	case algEs384:
		sig, err = esSign(findKey(keys, kid, ktyEc, useSig, opSign, alg), crypto.SHA384, headBodyPart)
	case algEs512:
		sig, err = esSign(findKey(keys, kid, ktyEc, useSig, opSign, alg), crypto.SHA512, headBodyPart)
	case algPs256:
		sig, err = psSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA256, headBodyPart)
	case algPs384:
		sig, err = psSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA384, headBodyPart)
	case algPs512:
		sig, err = psSign(findKey(keys, kid, ktyRsa, useSig, opSign, alg), crypto.SHA512, headBodyPart)
	default:
		return erro.New("unsupported sign algorithm " + alg)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	this.sig = sig
	this.headBodyPart = headBodyPart

	return nil
}

func (this *Jwt) Encrypt(keys []jwk.Key) error {
	if this.enced != nil {
		return nil
	}

	rawBody, err := this.getRawBody()
	if err != nil {
		return erro.Wrap(err)
	}

	alg, _ := this.Header(tagAlg).(string)
	enc, _ := this.Header(tagEnc).(string)

	var contKey []byte
	if alg != algDir {
		contKey, err = secrand.Bytes(keySizes[enc])
		if err != nil {
			return erro.Wrap(err)
		} else if len(contKey) == 0 {
			return erro.New("unsupported encryption " + enc)
		}
	}

	kid, _ := this.Header(tagKid).(string)

	var encedKey []byte
	switch alg {
	case algRsa1_5:
		encedKey, err = rsa15Encrypt(findKey(keys, kid, ktyRsa, useEnc, opWrapKey, alg), contKey)
	case algRsa_oaep:
		encedKey, err = rsaOaepEncrypt(findKey(keys, kid, ktyRsa, useEnc, opWrapKey, alg), crypto.SHA1, contKey)
	case algRsa_oaep_256:
		encedKey, err = rsaOaepEncrypt(findKey(keys, kid, ktyRsa, useEnc, opWrapKey, alg), crypto.SHA256, contKey)
	case algA128Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, ktyOct, useEnc, opWrapKey, alg), 16, contKey)
	case algA192Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, ktyOct, useEnc, opWrapKey, alg), 24, contKey)
	case algA256Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, ktyOct, useEnc, opWrapKey, alg), 32, contKey)
	case algDir:
		key := findKey(keys, kid, ktyOct, useEnc, opEncrypt)
		if key == nil {
			return erro.New("no key")
		} else if len(key.Common()) != keySizes[enc] {
			return erro.New("invalid key size")
		}
		contKey = key.Common()
		encedKey = []byte{}
	case algEcdh_es:
		encedKey, err = ecdhEsEncrypt(findKey(keys, kid, ktyEc, useEnc, opWrapKey, alg), contKey)
	case algEcdh_es_a128Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, ktyEc, useEnc, opWrapKey, alg), 16, contKey)
	case algEcdh_es_a192Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, ktyEc, useEnc, opWrapKey, alg), 24, contKey)
	case algEcdh_es_a256Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, ktyEc, useEnc, opWrapKey, alg), 32, contKey)
	case algA128Gcmkw, algA192Gcmkw, algA256Gcmkw:
		var initVec, authTag []byte
		initVec, encedKey, authTag, err = aGcmKwEncrypt(findKey(keys, kid, ktyOct, useEnc, opWrapKey, alg), keySizes[alg], contKey)
		if err == nil {
			this.SetHeader("iv", base64url.EncodeToString(initVec))  // 副作用注意。
			this.SetHeader("tag", base64url.EncodeToString(authTag)) // 副作用注意。
		}
	case algPbes2_hs256_a128Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", useEnc, opWrapKey, alg), crypto.SHA256, 16, contKey)
	case algPbes2_hs384_a192Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", useEnc, opWrapKey, alg), crypto.SHA384, 24, contKey)
	case algPbes2_hs512_a256Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", useEnc, opWrapKey, alg), crypto.SHA512, 32, contKey)
	default:
		return erro.New("unsupported key wrapping " + alg)
	}

	zip, _ := this.Header(tagZip).(string)

	var plain []byte
	switch zip {
	case "":
		plain = rawBody
	case zipDef:
		plain, err = defCompress(rawBody)
	default:
		return erro.New("Unsupported compression " + zip)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	rawHead, err := this.getRawHeader() // A128GCMKW の初期ベクトル等を追加するので、早めに実行しちゃ駄目。
	if err != nil {
		return erro.Wrap(err)
	}
	headPart := base64url.Encode(rawHead)

	var initVec, enced, authTag []byte
	switch enc {
	case encA128Cbc_Hs256:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 32, crypto.SHA256, plain, headPart)
	case encA192Cbc_Hs384:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 48, crypto.SHA384, plain, headPart)
	case encA256Cbc_Hs512:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 64, crypto.SHA512, plain, headPart)
	case encA128Gcm:
		initVec, enced, authTag, err = aGcmEncrypt(contKey, 16, plain, headPart)
	case encA192Gcm:
		initVec, enced, authTag, err = aGcmEncrypt(contKey, 24, plain, headPart)
	case encA256Gcm:
		initVec, enced, authTag, err = aGcmEncrypt(contKey, 32, plain, headPart)
	default:
		return erro.New("unsupported encryption" + enc)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	this.encedKey = encedKey
	this.initVec = initVec
	this.enced = enced
	this.authTag = authTag
	this.headPart = headPart

	return nil
}

// Compact serialization にする。
func (this *Jwt) Encode() ([]byte, error) {
	if this.sig != nil {
		// JWS.
		buff := append(this.headBodyPart, '.')
		this.compact = append(buff, base64url.Encode(this.sig)...)
	} else if this.enced != nil {
		// JWE.
		buff := append(this.headPart, '.')
		buff = append(buff, base64url.Encode(this.encedKey)...)
		buff = append(buff, '.')
		buff = append(buff, base64url.Encode(this.initVec)...)
		buff = append(buff, '.')
		buff = append(buff, base64url.Encode(this.enced)...)
		buff = append(buff, '.')
		this.compact = append(buff, base64url.Encode(this.authTag)...)
	} else {
		// 無署名 JWS
		rawHead, err := this.getRawHeader()
		if err != nil {
			return nil, erro.Wrap(err)
		}
		rawBody, err := this.getRawBody()
		if err != nil {
			return nil, erro.Wrap(err)
		}

		buff := append(base64url.Encode(rawHead), '.')
		buff = append(buff, base64url.Encode(rawBody)...)
		this.compact = append(buff, '.')
	}
	return this.compact, nil
}

// Compact serialization を読み取る。
func Parse(data []byte) (*Jwt, error) {
	switch parts := bytes.Split(data, []byte{'.'}); len(parts) {
	case 3:
		// JWS.
		rawHead, err := base64url.Decode(parts[0])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		rawBody, err := base64url.Decode(parts[1])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		sig, err := base64url.Decode(parts[2])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return &Jwt{
			rawHead:      rawHead,
			rawBody:      rawBody,
			sig:          sig,
			headBodyPart: data[:len(parts[0])+1+len(parts[1])],
			compact:      data,
		}, nil
	case 5:
		// JWE.
		rawHead, err := base64url.Decode(parts[0])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		encedKey, err := base64url.Decode(parts[1])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		initVec, err := base64url.Decode(parts[2])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		enced, err := base64url.Decode(parts[3])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		authTag, err := base64url.Decode(parts[4])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return &Jwt{
			rawHead:  rawHead,
			encedKey: encedKey,
			initVec:  initVec,
			enced:    enced,
			authTag:  authTag,
			headPart: parts[0],
			compact:  data,
		}, nil
	default:
		return nil, erro.New("invalid parts number")
	}
}

func (this *Jwt) RawHeader() []byte {
	if raw, err := this.getRawHeader(); err != nil {
		log.Warn(erro.Wrap(err))
		return []byte("{}")
	} else {
		return raw
	}
}

func (this *Jwt) getRawHeader() ([]byte, error) {
	if this.rawHead == nil {
		if this.head == nil {
			this.head = map[string]interface{}{}
		}
		var err error
		this.rawHead, err = json.Marshal(this.head)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return this.rawHead, nil
}

func (this *Jwt) RawBody() []byte {
	if raw, err := this.getRawBody(); err != nil {
		log.Warn(erro.Wrap(err))
		return []byte("{}")
	} else {
		return raw
	}
}

func (this *Jwt) getRawBody() ([]byte, error) {
	if this.rawBody == nil {
		if this.clms == nil {
			this.clms = map[string]interface{}{}
		}
		var err error
		this.rawBody, err = json.Marshal(this.clms)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return this.rawBody, nil
}

func (this *Jwt) IsSigned() bool {
	return this.sig != nil
}

func (this *Jwt) IsEncrypted() bool {
	return this.enced != nil
}

func (this *Jwt) Verify(keys []jwk.Key) (err error) {
	if this.sig == nil {
		return erro.New("not signed")
	}

	kid, _ := this.Header(tagKid).(string)
	alg, _ := this.Header(tagAlg).(string)

	switch alg {
	case algNone:
		err = noneVerify(this.sig)
	case algHs256:
		err = hsVerify(findKey(keys, kid, ktyOct, useSig, opVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case algHs384:
		err = hsVerify(findKey(keys, kid, ktyOct, useSig, opVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case algHs512:
		err = hsVerify(findKey(keys, kid, ktyOct, useSig, opVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case algRs256:
		err = rsVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case algRs384:
		err = rsVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case algRs512:
		err = rsVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case algEs256:
		err = esVerify(findKey(keys, kid, ktyEc, useSig, opVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case algEs384:
		err = esVerify(findKey(keys, kid, ktyEc, useSig, opVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case algEs512:
		err = esVerify(findKey(keys, kid, ktyEc, useSig, opVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case algPs256:
		err = psVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case algPs384:
		err = psVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case algPs512:
		err = psVerify(findKey(keys, kid, ktyRsa, useSig, opVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	default:
		return erro.New("unsupported sign " + alg)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	return nil
}

func (this *Jwt) Decrypt(keys []jwk.Key) (err error) {
	if this.enced == nil {
		return erro.New("not encrypted")
	}

	alg, _ := this.Header(tagAlg).(string)
	kid, _ := this.Header(tagKid).(string)

	var contKey []byte
	switch alg {
	case algRsa1_5:
		contKey, err = rsa15Decrypt(findKey(keys, kid, ktyRsa, useEnc, opUnwrapKey, alg), this.encedKey)
	case algRsa_oaep:
		contKey, err = rsaOaepDecrypt(findKey(keys, kid, ktyRsa, useEnc, opUnwrapKey, alg), crypto.SHA1, this.encedKey)
	case algRsa_oaep_256:
		contKey, err = rsaOaepDecrypt(findKey(keys, kid, ktyRsa, useEnc, opUnwrapKey, alg), crypto.SHA256, this.encedKey)
	case algA128Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, ktyOct, useEnc, opUnwrapKey, alg), 16, this.encedKey)
	case algA192Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, ktyOct, useEnc, opUnwrapKey, alg), 24, this.encedKey)
	case algA256Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, ktyOct, useEnc, opUnwrapKey, alg), 32, this.encedKey)
	case algDir:
		contKey, err = dirDecrypt(findKey(keys, kid, ktyOct, useEnc, opDecrypt, alg), this.encedKey)
	case algEcdh_es:
		contKey, err = ecdhEsDecrypt(findKey(keys, kid, ktyEc, useEnc, opUnwrapKey, alg), this.encedKey)
	case algEcdh_es_a128Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, ktyEc, useEnc, opUnwrapKey, alg), 16, this.encedKey)
	case algEcdh_es_a192Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, ktyEc, useEnc, opUnwrapKey, alg), 24, this.encedKey)
	case algEcdh_es_a256Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, ktyEc, useEnc, opUnwrapKey, alg), 32, this.encedKey)
	case algA128Gcmkw, algA192Gcmkw, algA256Gcmkw:
		var initVec, authTag []byte
		initVec, authTag, err = this.getInitVecAndAuthTagFromHeader()
		if err == nil {
			contKey, err = aGcmKwDecrypt(findKey(keys, kid, ktyOct, useEnc, opUnwrapKey, alg), keySizes[alg], initVec, this.encedKey, authTag)
		}
	case algPbes2_hs256_a128Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", useEnc, opUnwrapKey, alg), crypto.SHA256, 16, this.encedKey)
	case algPbes2_hs384_a192Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", useEnc, opUnwrapKey, alg), crypto.SHA384, 24, this.encedKey)
	case algPbes2_hs512_a256Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", useEnc, opUnwrapKey, alg), crypto.SHA512, 32, this.encedKey)
	default:
		return erro.New("unsupported key wrapping " + alg)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	enc, _ := this.Header(tagEnc).(string)

	var plain []byte
	switch enc {
	case encA128Cbc_Hs256:
		plain, err = aCbcHsDecrypt(contKey, 32, crypto.SHA256, this.headPart, this.initVec, this.enced, this.authTag)
	case encA192Cbc_Hs384:
		plain, err = aCbcHsDecrypt(contKey, 48, crypto.SHA384, this.headPart, this.initVec, this.enced, this.authTag)
	case encA256Cbc_Hs512:
		plain, err = aCbcHsDecrypt(contKey, 64, crypto.SHA512, this.headPart, this.initVec, this.enced, this.authTag)
	case encA128Gcm:
		plain, err = aGcmDecrypt(contKey, 16, this.headPart, this.initVec, this.enced, this.authTag)
	case encA192Gcm:
		plain, err = aGcmDecrypt(contKey, 24, this.headPart, this.initVec, this.enced, this.authTag)
	case encA256Gcm:
		plain, err = aGcmDecrypt(contKey, 32, this.headPart, this.initVec, this.enced, this.authTag)
	default:
		return erro.New("unsupported encryption " + enc)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	zip, _ := this.Header(tagZip).(string)

	var rawBody []byte
	switch zip {
	case "":
		rawBody = plain
	case "DEF":
		rawBody, err = defDecompress(plain)
	default:
		return erro.New("unsupported compression " + zip)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	this.rawBody = rawBody

	return nil
}

func (this *Jwt) Header(tag string) interface{} {
	if val, err := this.getHeader(tag); err != nil {
		log.Warn(erro.Unwrap(err))
		log.Debug(erro.Wrap(err))
		return nil
	} else {
		return val
	}
}

func (this *Jwt) getHeader(tag string) (interface{}, error) {
	if this.head == nil {
		var err error
		this.head, err = parseJson(this.rawHead)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return this.head[tag], nil
}

func (this *Jwt) Claim(tag string) interface{} {
	if val, err := this.getClaim(tag); err != nil {
		log.Warn(erro.Unwrap(err))
		log.Debug(erro.Wrap(err))
		return nil
	} else {
		return val
	}
}

func (this *Jwt) getClaim(tag string) (interface{}, error) {
	if this.clms == nil {
		var err error
		this.clms, err = parseJson(this.rawBody)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return this.clms[tag], nil
}

func parseJson(data []byte) (map[string]interface{}, error) {
	if data == nil {
		return map[string]interface{}{}, nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, erro.Wrap(err)
	}
	return m, nil
}

func parseJsonOrNew(data []byte) map[string]interface{} {
	if m, err := parseJson(data); err != nil {
		log.Warn(erro.Unwrap(err))
		log.Debug(erro.Wrap(err))
		return map[string]interface{}{}
	} else {
		return m
	}
}

func (this *Jwt) getInitVecAndAuthTagFromHeader() (initVec, authTag []byte, err error) {
	if str, _ := this.Header(tagIv).(string); str == "" {
		return nil, nil, erro.New("no initialization vector")
	} else if initVec, err := base64url.DecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if str, _ := this.Header(tagTag).(string); str == "" {
		return nil, nil, erro.New("no authentication tag")
	} else if authTag, err := base64url.DecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	} else {
		return initVec, authTag, nil
	}
}
