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

// JWT 関係。
package jwt

import (
	"bytes"
	"crypto"
	_ "crypto/ecdsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/json"

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

	headBodyPart := base64UrlEncode(rawHead)
	headBodyPart = append(headBodyPart, '.')
	headBodyPart = append(headBodyPart, base64UrlEncode(rawBody)...)

	kid, _ := this.Header(tagKid).(string)
	alg, _ := this.Header(tagAlg).(string)

	var sig []byte
	switch alg {
	case tagNone:
		sig = []byte{}
	case tagHs256:
		sig, err = hsSign(findKey(keys, kid, tagOct, tagSig, tagSign, alg), crypto.SHA256, headBodyPart)
	case tagHs384:
		sig, err = hsSign(findKey(keys, kid, tagOct, tagSig, tagSign, alg), crypto.SHA384, headBodyPart)
	case tagHs512:
		sig, err = hsSign(findKey(keys, kid, tagOct, tagSig, tagSign, alg), crypto.SHA512, headBodyPart)
	case tagRs256:
		sig, err = rsSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA256, headBodyPart)
	case tagRs384:
		sig, err = rsSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA384, headBodyPart)
	case tagRs512:
		sig, err = rsSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA512, headBodyPart)
	case tagEs256:
		sig, err = esSign(findKey(keys, kid, tagEc, tagSig, tagSign, alg), crypto.SHA256, headBodyPart)
	case tagEs384:
		sig, err = esSign(findKey(keys, kid, tagEc, tagSig, tagSign, alg), crypto.SHA384, headBodyPart)
	case tagEs512:
		sig, err = esSign(findKey(keys, kid, tagEc, tagSig, tagSign, alg), crypto.SHA512, headBodyPart)
	case tagPs256:
		sig, err = psSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA256, headBodyPart)
	case tagPs384:
		sig, err = psSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA384, headBodyPart)
	case tagPs512:
		sig, err = psSign(findKey(keys, kid, tagRsa, tagSig, tagSign, alg), crypto.SHA512, headBodyPart)
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
	if alg != tagDir {
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
	case tagRsa1_5:
		encedKey, err = rsa15Encrypt(findKey(keys, kid, tagRsa, tagEnc, tagWrapKey, alg), contKey)
	case tagRsa_oaep:
		encedKey, err = rsaOaepEncrypt(findKey(keys, kid, tagRsa, tagEnc, tagWrapKey, alg), crypto.SHA1, contKey)
	case tagRsa_oaep_256:
		encedKey, err = rsaOaepEncrypt(findKey(keys, kid, tagRsa, tagEnc, tagWrapKey, alg), crypto.SHA256, contKey)
	case tagA128Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, tagOct, tagEnc, tagWrapKey, alg), 16, contKey)
	case tagA192Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, tagOct, tagEnc, tagWrapKey, alg), 24, contKey)
	case tagA256Kw:
		encedKey, err = aKwEncrypt(findKey(keys, kid, tagOct, tagEnc, tagWrapKey, alg), 32, contKey)
	case tagDir:
		key := findKey(keys, kid, tagOct, tagEnc, tagEncrypt)
		if key == nil {
			return erro.New("no key")
		} else if len(key.Common()) != keySizes[enc] {
			return erro.New("invalid key size")
		}
		contKey = key.Common()
		encedKey = []byte{}
	case tagEcdh_es:
		encedKey, err = ecdhEsEncrypt(findKey(keys, kid, tagEc, tagEnc, tagWrapKey, alg), contKey)
	case tagEcdh_es_a128Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, tagEc, tagEnc, tagWrapKey, alg), 16, contKey)
	case tagEcdh_es_a192Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, tagEc, tagEnc, tagWrapKey, alg), 24, contKey)
	case tagEcdh_es_a256Kw:
		encedKey, err = ecdhEsAKwEncrypt(findKey(keys, kid, tagEc, tagEnc, tagWrapKey, alg), 32, contKey)
	case tagA128Gcmkw, tagA192Gcmkw, tagA256Gcmkw:
		var initVec, authTag []byte
		initVec, encedKey, authTag, err = aGcmKwEncrypt(findKey(keys, kid, tagOct, tagEnc, tagWrapKey, alg), keySizes[alg], contKey)
		if err == nil {
			this.SetHeader("iv", base64.RawURLEncoding.EncodeToString(initVec))  // 副作用注意。
			this.SetHeader("tag", base64.RawURLEncoding.EncodeToString(authTag)) // 副作用注意。
		}
	case tagPbes2_hs256_a128Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", tagEnc, tagWrapKey, alg), crypto.SHA256, 16, contKey)
	case tagPbes2_hs384_a192Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", tagEnc, tagWrapKey, alg), crypto.SHA384, 24, contKey)
	case tagPbes2_hs512_a256Kw:
		encedKey, err = pbes2HsAKwEncrypt(findKey(keys, kid, "", tagEnc, tagWrapKey, alg), crypto.SHA512, 32, contKey)
	default:
		return erro.New("unsupported key wrapping " + alg)
	}

	zip, _ := this.Header(tagZip).(string)

	var plain []byte
	switch zip {
	case "":
		plain = rawBody
	case tagDef:
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
	headPart := base64UrlEncode(rawHead)

	var initVec, enced, authTag []byte
	switch enc {
	case tagA128Cbc_Hs256:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 32, crypto.SHA256, plain, headPart)
	case tagA192Cbc_Hs384:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 48, crypto.SHA384, plain, headPart)
	case tagA256Cbc_Hs512:
		initVec, enced, authTag, err = aCbcHsEncrypt(contKey, 64, crypto.SHA512, plain, headPart)
	case tagA128Gcm:
		initVec, enced, authTag, err = aGcmEncrypt(contKey, 16, plain, headPart)
	case tagA192Gcm:
		initVec, enced, authTag, err = aGcmEncrypt(contKey, 24, plain, headPart)
	case tagA256Gcm:
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
		this.compact = append(buff, base64UrlEncode(this.sig)...)
	} else if this.enced != nil {
		// JWE.
		buff := append(this.headPart, '.')
		buff = append(buff, base64UrlEncode(this.encedKey)...)
		buff = append(buff, '.')
		buff = append(buff, base64UrlEncode(this.initVec)...)
		buff = append(buff, '.')
		buff = append(buff, base64UrlEncode(this.enced)...)
		buff = append(buff, '.')
		this.compact = append(buff, base64UrlEncode(this.authTag)...)
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

		buff := append(base64UrlEncode(rawHead), '.')
		buff = append(buff, base64UrlEncode(rawBody)...)
		this.compact = append(buff, '.')
	}
	return this.compact, nil
}

// Compact serialization を読み取る。
func Parse(data []byte) (*Jwt, error) {
	switch parts := bytes.Split(data, []byte{'.'}); len(parts) {
	case 3:
		// JWS.
		rawHead, err := base64UrlDecode(parts[0])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		rawBody, err := base64UrlDecode(parts[1])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		sig, err := base64UrlDecode(parts[2])
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
		rawHead, err := base64UrlDecode(parts[0])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		encedKey, err := base64UrlDecode(parts[1])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		initVec, err := base64UrlDecode(parts[2])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		enced, err := base64UrlDecode(parts[3])
		if err != nil {
			return nil, erro.Wrap(err)
		}
		authTag, err := base64UrlDecode(parts[4])
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
	case tagNone:
		err = noneVerify(this.sig)
	case tagHs256:
		err = hsVerify(findKey(keys, kid, tagOct, tagSig, tagVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case tagHs384:
		err = hsVerify(findKey(keys, kid, tagOct, tagSig, tagVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case tagHs512:
		err = hsVerify(findKey(keys, kid, tagOct, tagSig, tagVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case tagRs256:
		err = rsVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case tagRs384:
		err = rsVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case tagRs512:
		err = rsVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case tagEs256:
		err = esVerify(findKey(keys, kid, tagEc, tagSig, tagVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case tagEs384:
		err = esVerify(findKey(keys, kid, tagEc, tagSig, tagVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case tagEs512:
		err = esVerify(findKey(keys, kid, tagEc, tagSig, tagVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
	case tagPs256:
		err = psVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA256, this.sig, this.headBodyPart)
	case tagPs384:
		err = psVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA384, this.sig, this.headBodyPart)
	case tagPs512:
		err = psVerify(findKey(keys, kid, tagRsa, tagSig, tagVerify, alg), crypto.SHA512, this.sig, this.headBodyPart)
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
	case tagRsa1_5:
		contKey, err = rsa15Decrypt(findKey(keys, kid, tagRsa, tagEnc, tagUnwrapKey, alg), this.encedKey)
	case tagRsa_oaep:
		contKey, err = rsaOaepDecrypt(findKey(keys, kid, tagRsa, tagEnc, tagUnwrapKey, alg), crypto.SHA1, this.encedKey)
	case tagRsa_oaep_256:
		contKey, err = rsaOaepDecrypt(findKey(keys, kid, tagRsa, tagEnc, tagUnwrapKey, alg), crypto.SHA256, this.encedKey)
	case tagA128Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, tagOct, tagEnc, tagUnwrapKey, alg), 16, this.encedKey)
	case tagA192Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, tagOct, tagEnc, tagUnwrapKey, alg), 24, this.encedKey)
	case tagA256Kw:
		contKey, err = aKwDecrypt(findKey(keys, kid, tagOct, tagEnc, tagUnwrapKey, alg), 32, this.encedKey)
	case tagDir:
		contKey, err = dirDecrypt(findKey(keys, kid, tagOct, tagEnc, tagDecrypt, alg), this.encedKey)
	case tagEcdh_es:
		contKey, err = ecdhEsDecrypt(findKey(keys, kid, tagEc, tagEnc, tagUnwrapKey, alg), this.encedKey)
	case tagEcdh_es_a128Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, tagEc, tagEnc, tagUnwrapKey, alg), 16, this.encedKey)
	case tagEcdh_es_a192Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, tagEc, tagEnc, tagUnwrapKey, alg), 24, this.encedKey)
	case tagEcdh_es_a256Kw:
		contKey, err = ecdhEsAKwDecrypt(findKey(keys, kid, tagEc, tagEnc, tagUnwrapKey, alg), 32, this.encedKey)
	case tagA128Gcmkw, tagA192Gcmkw, tagA256Gcmkw:
		var initVec, authTag []byte
		initVec, authTag, err = this.getInitVecAndAuthTagFromHeader()
		if err == nil {
			contKey, err = aGcmKwDecrypt(findKey(keys, kid, tagOct, tagEnc, tagUnwrapKey, alg), keySizes[alg], initVec, this.encedKey, authTag)
		}
	case tagPbes2_hs256_a128Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", tagEnc, tagUnwrapKey, alg), crypto.SHA256, 16, this.encedKey)
	case tagPbes2_hs384_a192Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", tagEnc, tagUnwrapKey, alg), crypto.SHA384, 24, this.encedKey)
	case tagPbes2_hs512_a256Kw:
		contKey, err = pbes2HsAKwDecrypt(findKey(keys, kid, "", tagEnc, tagUnwrapKey, alg), crypto.SHA512, 32, this.encedKey)
	default:
		return erro.New("unsupported key wrapping " + alg)
	}
	if err != nil {
		return erro.Wrap(err)
	}

	enc, _ := this.Header(tagEnc).(string)

	var plain []byte
	switch enc {
	case tagA128Cbc_Hs256:
		plain, err = aCbcHsDecrypt(contKey, 32, crypto.SHA256, this.headPart, this.initVec, this.enced, this.authTag)
	case tagA192Cbc_Hs384:
		plain, err = aCbcHsDecrypt(contKey, 48, crypto.SHA384, this.headPart, this.initVec, this.enced, this.authTag)
	case tagA256Cbc_Hs512:
		plain, err = aCbcHsDecrypt(contKey, 64, crypto.SHA512, this.headPart, this.initVec, this.enced, this.authTag)
	case tagA128Gcm:
		plain, err = aGcmDecrypt(contKey, 16, this.headPart, this.initVec, this.enced, this.authTag)
	case tagA192Gcm:
		plain, err = aGcmDecrypt(contKey, 24, this.headPart, this.initVec, this.enced, this.authTag)
	case tagA256Gcm:
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

func (this *Jwt) ClaimNames() []string {
	if this.clms == nil {
		var err error
		this.clms, err = parseJson(this.rawBody)
		if err != nil {
			log.Warn(erro.Unwrap(err))
			log.Debug(erro.Wrap(err))
		}
	}
	names := []string{}
	for name := range this.clms {
		names = append(names, name)
	}
	return names
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
	} else if initVec, err := base64.RawURLEncoding.DecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if str, _ := this.Header(tagTag).(string); str == "" {
		return nil, nil, erro.New("no authentication tag")
	} else if authTag, err := base64.RawURLEncoding.DecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	} else {
		return initVec, authTag, nil
	}
}

func base64UrlEncode(src []byte) []byte {
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(dst, src)
	return dst
}

func base64UrlDecode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(src)))
	n, err := base64.RawURLEncoding.Decode(dst, src)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return dst[:n], nil
}
