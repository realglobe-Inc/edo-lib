package jwt

import (
	"crypto"
	"crypto/ecdsa"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/realglobe-Inc/edo-toolkit/util/secrand"
	"github.com/realglobe-Inc/go-lib/erro"
	"strings"
)

// JSON Web Token
type Jwt struct {
	head   map[string]interface{}
	clms   map[string]interface{}
	nested *Jwt

	encoded []byte
}

func New() *Jwt {
	return &Jwt{
		head: map[string]interface{}{},
		clms: map[string]interface{}{},
	}
}

func (this *Jwt) HeaderNames() map[string]bool {
	m := map[string]bool{}
	for k := range this.head {
		m[k] = true
	}
	return m
}

func (this *Jwt) Header(headName string) interface{} {
	return this.head[headName]
}

// val が nil や空文字列の場合は削除する。
func (this *Jwt) SetHeader(tag string, val interface{}) {
	this.encoded = nil
	if val == nil || val == "" {
		delete(this.head, tag)
	} else {
		this.head[tag] = val
	}
}

func (this *Jwt) ClaimNames() map[string]bool {
	m := map[string]bool{}
	for k := range this.clms {
		m[k] = true
	}
	return m
}

// クレームを返す。Nesting() が true のときは必ず nil を返す。
func (this *Jwt) Claim(clmName string) interface{} {
	return this.clms[clmName]
}

// val が nil や空文字列の場合は削除する。
func (this *Jwt) SetClaim(tag string, val interface{}) {
	this.nested = nil
	this.encoded = nil
	if val == nil || val == "" {
		delete(this.clms, tag)
	} else {
		this.clms[tag] = val
	}
}

// 入れ子かどうか。
func (this *Jwt) Nesting() bool {
	return this.nested != nil
}

// 入れ子にされた Jwt を返す。Nesting() が true のときのみ非 nil を返す。
func (this *Jwt) Nested() *Jwt {
	return this.nested
}

// 入れ子にする。入れ子にしたあとで jt を操作したときの動作は定義しない。
func (this *Jwt) Nest(jt *Jwt) {
	if jt == this.nested {
		return
	}
	if len(this.clms) > 0 {
		this.clms = map[string]interface{}{}
	}
	this.encoded = nil
	this.nested = jt
}

// Compact serialization.
func (this *Jwt) Encode(sigKeys, encKeys map[string]interface{}) ([]byte, error) {
	if this.encoded != nil {
		return this.encoded, nil
	} else if alg, _ := this.Header("alg").(string); isJwsAlgorithm(alg) {
		return this.jwsEncode(alg, sigKeys)
	} else if isJweAlgorithm(alg) {
		return this.jweEncode(alg, sigKeys, encKeys)
	} else {
		return nil, erro.New("alg " + alg + " is not unsupported")
	}
}

func (this *Jwt) jwsEncode(alg string, sigKeys map[string]interface{}) ([]byte, error) {

	var key interface{}
	if alg != "none" {
		kid, _ := this.Header("kid").(string)
		var err error
		key, err = _getKey(kid, sigKeys)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	var buff []byte
	if data, err := json.Marshal(this.head); err != nil {
		return nil, erro.Wrap(err)
	} else {
		buff = base64UrlEncode(data)
	}
	buff = append(buff, '.')
	if data, err := json.Marshal(this.clms); err != nil {
		return nil, erro.Wrap(err)
	} else {
		buff = append(buff, base64UrlEncode(data)...)
	}

	var sig []byte
	var err error
	switch alg {
	case "none":
	case "HS256":
		sig, err = hsSign(key, crypto.SHA256, buff)
	case "HS384":
		sig, err = hsSign(key, crypto.SHA384, buff)
	case "HS512":
		sig, err = hsSign(key, crypto.SHA512, buff)
	case "RS256":
		sig, err = rsSign(key, crypto.SHA256, buff)
	case "RS384":
		sig, err = rsSign(key, crypto.SHA384, buff)
	case "RS512":
		sig, err = rsSign(key, crypto.SHA512, buff)
	case "ES256":
		// JWA の仕様で ESxxx は鍵のサイズが決められている。
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 256 {
			return nil, erro.New("not P-256 EC key")
		}
		sig, err = esSign(priKey, crypto.SHA256, buff)
	case "ES384":
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 384 {
			return nil, erro.New("not P-384 EC key")
		}
		sig, err = esSign(priKey, crypto.SHA384, buff)
	case "ES512":
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 521 {
			return nil, erro.New("not P-521 EC key")
		}
		sig, err = esSign(priKey, crypto.SHA512, buff)
	case "PS256":
		sig, err = psSign(key, crypto.SHA256, buff)
	case "PS384":
		sig, err = psSign(key, crypto.SHA384, buff)
	case "PS512":
		sig, err = psSign(key, crypto.SHA512, buff)
	default:
		return nil, erro.New("alg " + alg + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	buff = append(buff, '.')
	this.encoded = append(buff, base64UrlEncode(sig)...)

	return this.encoded, nil
}

func (this *Jwt) jweEncode(alg string, sigKeys, encKeys map[string]interface{}) ([]byte, error) {
	kid, _ := this.Header("kid").(string)
	key, err := _getKey(kid, encKeys)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var plain []byte
	if this.Nesting() {
		this.SetHeader("cty", "JWT") // 副作用注意。
		plain, err = this.Nested().Encode(sigKeys, encKeys)
	} else {
		plain, err = json.Marshal(this.clms)
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	switch zip, _ := this.Header("zip").(string); zip {
	case "":
	case "DEF":
		plain, err = defCompress(plain)
	default:
		return nil, erro.New("zip " + zip + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	enc, _ := this.Header("enc").(string)

	var contKey []byte
	if alg == "dir" {
		var ok bool
		contKey, ok = key.([]byte)
		if !ok {
			return nil, erro.New("key cannot be used for dir")
		}

		switch enc {
		case "A128CBC-HS256":
			ok = (len(contKey) == 32)
		case "A192CBC-HS384":
			ok = (len(contKey) == 48)
		case "A256CBC-HS512":
			ok = (len(contKey) == 64)
		case "A128GCM":
			ok = (len(contKey) == 16)
		case "A192GCM":
			ok = (len(contKey) == 24)
		case "A256GCM":
			ok = (len(contKey) == 32)
		default:
			return nil, erro.New("enc " + enc + " is unsupported")
		}
		if !ok {
			return nil, erro.New("invalid key size")
		}
	} else {
		switch enc {
		case "A128CBC-HS256":
			contKey, err = secrand.Bytes(32)
		case "A192CBC-HS384":
			contKey, err = secrand.Bytes(48)
		case "A256CBC-HS512":
			contKey, err = secrand.Bytes(64)
		case "A128GCM":
			contKey, err = secrand.Bytes(16)
		case "A192GCM":
			contKey, err = secrand.Bytes(24)
		case "A256GCM":
			contKey, err = secrand.Bytes(32)
		default:
			return nil, erro.New("enc " + enc + " is unsupported")
		}
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	var encryptedKey []byte
	switch alg {
	case "RSA1_5":
		encryptedKey, err = rsa15Encrypt(key, contKey)
	case "RSA-OAEP":
		encryptedKey, err = rsaOaepEncrypt(key, crypto.SHA1, contKey)
	case "RSA-OAEP-256":
		encryptedKey, err = rsaOaepEncrypt(key, crypto.SHA256, contKey)
	case "A128KW":
		encryptedKey, err = aKwEncrypt(key, 16, contKey)
	case "A192KW":
		encryptedKey, err = aKwEncrypt(key, 24, contKey)
	case "A256KW":
		encryptedKey, err = aKwEncrypt(key, 32, contKey)
	case "dir":
		encryptedKey, err = dirEncrypt(key)
	case "ECDH-ES":
		encryptedKey, err = ecdhEsEncrypt(key, contKey)
	case "ECDH-ES+A128KW":
		encryptedKey, err = ecdhEsAKwEncrypt(key, 16, contKey)
	case "ECDH-ES+A192KW":
		encryptedKey, err = ecdhEsAKwEncrypt(key, 24, contKey)
	case "ECDH-ES+A256KW":
		encryptedKey, err = ecdhEsAKwEncrypt(key, 32, contKey)
	case "A128GCMKW":
		var initVec, authTag []byte
		initVec, encryptedKey, authTag, err = aGcmKwEncrypt(key, 16, contKey)
		if err == nil {
			this.SetHeader("iv", base64UrlEncodeToString(initVec))  // 副作用注意。
			this.SetHeader("tag", base64UrlEncodeToString(authTag)) // 副作用注意。
		}
	case "A192GCMKW":
		var initVec, authTag []byte
		initVec, encryptedKey, authTag, err = aGcmKwEncrypt(key, 24, contKey)
		if err == nil {
			this.SetHeader("iv", base64UrlEncodeToString(initVec))  // 副作用注意。
			this.SetHeader("tag", base64UrlEncodeToString(authTag)) // 副作用注意。
		}
	case "A256GCMKW":
		var initVec, authTag []byte
		initVec, encryptedKey, authTag, err = aGcmKwEncrypt(key, 32, contKey)
		if err == nil {
			this.SetHeader("iv", base64UrlEncodeToString(initVec))  // 副作用注意。
			this.SetHeader("tag", base64UrlEncodeToString(authTag)) // 副作用注意。
		}
	case "PBES2-HS256+A128KW":
		encryptedKey, err = pbes2HsAKwEncrypt(key, crypto.SHA256, 16, contKey)
	case "PBES2-HS384+A192KW":
		encryptedKey, err = pbes2HsAKwEncrypt(key, crypto.SHA384, 24, contKey)
	case "PBES2-HS512+A256KW":
		encryptedKey, err = pbes2HsAKwEncrypt(key, crypto.SHA512, 32, contKey)
	default:
		return nil, erro.New("alg " + alg + " is unsupported")
	}

	var headPart []byte
	if data, err := json.Marshal(this.head); err != nil {
		return nil, erro.Wrap(err)
	} else {
		headPart = base64UrlEncode(data)
	}

	var initVec, encrypted, authTag []byte
	switch enc {
	case "A128CBC-HS256":
		initVec, encrypted, authTag, err = aCbcHsEncrypt(contKey, 32, crypto.SHA256, plain, headPart)
	case "A192CBC-HS384":
		initVec, encrypted, authTag, err = aCbcHsEncrypt(contKey, 48, crypto.SHA384, plain, headPart)
	case "A256CBC-HS512":
		initVec, encrypted, authTag, err = aCbcHsEncrypt(contKey, 64, crypto.SHA512, plain, headPart)
	case "A128GCM":
		initVec, encrypted, authTag, err = aGcmEncrypt(contKey, 16, plain, headPart)
	case "A192GCM":
		initVec, encrypted, authTag, err = aGcmEncrypt(contKey, 24, plain, headPart)
	case "A256GCM":
		initVec, encrypted, authTag, err = aGcmEncrypt(contKey, 32, plain, headPart)
	default:
		return nil, erro.New("enc " + enc + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	buff := append(headPart, '.')
	buff = append(buff, base64UrlEncode(encryptedKey)...)
	buff = append(buff, '.')
	buff = append(buff, base64UrlEncode(initVec)...)
	buff = append(buff, '.')
	buff = append(buff, base64UrlEncode(encrypted)...)
	buff = append(buff, '.')
	this.encoded = append(buff, base64UrlEncode(authTag)...)

	return this.encoded, nil
}

// JSON serialization.
func (this *Jwt) EncodeToJson() ([]byte, error) {
	panic("not yet implemented")
}

func Parse(encoded string, veriKeys, decKeys map[string]interface{}) (*Jwt, error) {
	var jt *Jwt
	switch parts := strings.Split(encoded, "."); len(parts) {
	case 3:
		var err error
		jt, err = parseJwsParts(parts[0], parts[1], parts[2], veriKeys)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	case 5:
		var err error
		jt, err = parseJweParts(parts[0], parts[1], parts[2], parts[3], parts[4], veriKeys, decKeys)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	default:
		return nil, erro.New("invalid parts number")
	}
	jt.encoded = []byte(encoded)
	return jt, nil
}

func parseJwsParts(headPart, clmsPart, sigPart string, veriKeys map[string]interface{}) (*Jwt, error) {
	var head map[string]interface{}
	if data, err := base64UrlDecodeString(headPart); err != nil {
		return nil, erro.Wrap(err)
	} else if err := json.Unmarshal(data, &head); err != nil {
		return nil, erro.Wrap(err)
	}

	var clms map[string]interface{}
	if data, err := base64UrlDecodeString(clmsPart); err != nil {
		return nil, erro.Wrap(err)
	} else if err := json.Unmarshal(data, &clms); err != nil {
		return nil, erro.Wrap(err)
	}

	alg, _ := head["alg"].(string)

	var key interface{}
	if alg != "none" {
		kid, _ := head["kid"].(string)
		var err error
		key, err = _getKey(kid, veriKeys)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	sig, err := base64UrlDecodeString(sigPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	switch alg {
	case "none":
		err = noneVerify(sig)
	case "HS256":
		err = hsVerify(key, crypto.SHA256, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "HS384":
		err = hsVerify(key, crypto.SHA384, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "HS512":
		err = hsVerify(key, crypto.SHA512, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "RS256":
		err = rsVerify(key, crypto.SHA256, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "RS384":
		err = rsVerify(key, crypto.SHA384, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "RS512":
		err = rsVerify(key, crypto.SHA512, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "ES256":
		// JWA の仕様で ESxxx は鍵のサイズが決められている。
		priKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 256 {
			return nil, erro.New("not P-256 EC key")
		}
		err = esVerify(priKey, crypto.SHA256, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "ES384":
		priKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 384 {
			return nil, erro.New("not P-384 EC key")
		}
		err = esVerify(priKey, crypto.SHA384, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "ES512":
		priKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 521 {
			return nil, erro.New("not P-521 EC key")
		}
		err = esVerify(priKey, crypto.SHA512, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "PS256":
		err = psVerify(key, crypto.SHA256, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "PS384":
		err = psVerify(key, crypto.SHA384, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	case "PS512":
		err = psVerify(key, crypto.SHA512, sig, []byte(headPart), []byte{'.'}, []byte(clmsPart))
	default:
		return nil, erro.New("alg " + alg + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return &Jwt{head, clms, nil, nil}, nil
}

func parseJweParts(headPart, encryptedKeyPart, initVecPart, encryptedPart, authTagPart string, veriKeys, decKeys map[string]interface{}) (*Jwt, error) {
	var head map[string]interface{}
	if data, err := base64UrlDecodeString(headPart); err != nil {
		return nil, erro.Wrap(err)
	} else if err := json.Unmarshal(data, &head); err != nil {
		return nil, erro.Wrap(err)
	}

	encryptedKey, err := base64UrlDecodeString(encryptedKeyPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	initVec, err := base64UrlDecodeString(initVecPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	encrypted, err := base64UrlDecodeString(encryptedPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	authTag, err := base64UrlDecodeString(authTagPart)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	alg, _ := head["alg"].(string)

	kid, _ := head["kid"].(string)
	key, err := _getKey(kid, decKeys)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var contKey []byte
	switch alg {
	case "RSA1_5":
		contKey, err = rsa15Decrypt(key, encryptedKey)
	case "RSA-OAEP":
		contKey, err = rsaOaepDecrypt(key, crypto.SHA1, encryptedKey)
	case "RSA-OAEP-256":
		contKey, err = rsaOaepDecrypt(key, crypto.SHA256, encryptedKey)
	case "A128KW":
		contKey, err = aKwDecrypt(key, 16, encryptedKey)
	case "A192KW":
		contKey, err = aKwDecrypt(key, 24, encryptedKey)
	case "A256KW":
		contKey, err = aKwDecrypt(key, 32, encryptedKey)
	case "dir":
		contKey, err = dirDecrypt(key, encryptedKey)
	case "ECDH-ES":
		contKey, err = ecdhEsDecrypt(key, encryptedKey)
	case "ECDH-ES+A128KW":
		contKey, err = ecdhEsAKwDecrypt(key, 16, encryptedKey)
	case "ECDH-ES+A192KW":
		contKey, err = ecdhEsAKwDecrypt(key, 24, encryptedKey)
	case "ECDH-ES+A256KW":
		contKey, err = ecdhEsAKwDecrypt(key, 32, encryptedKey)
	case "A128GCMKW":
		var initVec, authTag []byte
		initVec, authTag, err = _getInitVecAndAuthTag(head)
		if err == nil {
			contKey, err = aGcmKwDecrypt(key, 16, initVec, encryptedKey, authTag)
		}
	case "A192GCMKW":
		var initVec, authTag []byte
		initVec, authTag, err = _getInitVecAndAuthTag(head)
		if err == nil {
			contKey, err = aGcmKwDecrypt(key, 24, initVec, encryptedKey, authTag)
		}
	case "A256GCMKW":
		var initVec, authTag []byte
		initVec, authTag, err = _getInitVecAndAuthTag(head)
		if err == nil {
			contKey, err = aGcmKwDecrypt(key, 32, initVec, encryptedKey, authTag)
		}
	case "PBES2-HS256+A128KW":
		contKey, err = pbes2HsAKwDecrypt(key, crypto.SHA256, 16, encryptedKey)
	case "PBES2-HS384+A192KW":
		contKey, err = pbes2HsAKwDecrypt(key, crypto.SHA384, 24, encryptedKey)
	case "PBES2-HS512+A256KW":
		contKey, err = pbes2HsAKwDecrypt(key, crypto.SHA512, 32, encryptedKey)
	default:
		return nil, erro.New("alg " + alg + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var plain []byte
	switch enc, _ := head["enc"].(string); enc {
	case "A128CBC-HS256":
		plain, err = aCbcHsDecrypt(contKey, 32, crypto.SHA256, []byte(headPart), initVec, encrypted, authTag)
	case "A192CBC-HS384":
		plain, err = aCbcHsDecrypt(contKey, 48, crypto.SHA384, []byte(headPart), initVec, encrypted, authTag)
	case "A256CBC-HS512":
		plain, err = aCbcHsDecrypt(contKey, 64, crypto.SHA512, []byte(headPart), initVec, encrypted, authTag)
	case "A128GCM":
		plain, err = aGcmDecrypt(contKey, 16, []byte(headPart), initVec, encrypted, authTag)
	case "A192GCM":
		plain, err = aGcmDecrypt(contKey, 24, []byte(headPart), initVec, encrypted, authTag)
	case "A256GCM":
		plain, err = aGcmDecrypt(contKey, 32, []byte(headPart), initVec, encrypted, authTag)
	default:
		return nil, erro.New("enc " + enc + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	switch zip, _ := head["zip"].(string); zip {
	case "":
	case "DEF":
		plain, err = defDecompress(plain)
	default:
		return nil, erro.New("zip " + zip + " is unsupported")
	}
	if err != nil {
		return nil, erro.Wrap(err)
	}

	switch cty, _ := head["cty"].(string); cty {
	case "":
		var clms map[string]interface{}
		if err := json.Unmarshal(plain, &clms); err != nil {
			fmt.Println("Aho ", string(plain))
			return nil, erro.Wrap(err)
		}
		return &Jwt{head, clms, nil, nil}, nil
	case "JWT":
		jt, err := Parse(string(plain), veriKeys, decKeys)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return &Jwt{head, map[string]interface{}{}, jt, nil}, nil
	default:
		return nil, erro.New("cty " + cty + " is unsupported")
	}
}

func _getKey(kid string, keys map[string]interface{}) (interface{}, error) {
	if key := keys[kid]; key != nil {
		return key, nil
	}

	if kid != "" {
		return nil, erro.New("no key " + kid)
	} else if len(keys) == 1 {
		// 1 つだけならそれを使う。
		for _, key := range keys {
			return key, nil
		}
	}
	return nil, erro.New("key is not specified")
}

func _getInitVecAndAuthTag(head map[string]interface{}) (initVec, authTag []byte, err error) {
	if str, _ := head["iv"].(string); str == "" {
		return nil, nil, erro.New("no iv")
	} else if initVec, err = base64UrlDecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if str, _ := head["tag"].(string); str == "" {
		return nil, nil, erro.New("no tag")
	} else if authTag, err = base64UrlDecodeString(str); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return initVec, authTag, nil
}

// ヘッダとクレーム部を JSON 形式で返す。
// 入れ子ならヘッダは複数になる
func ToJsons(jt *Jwt) ([][]byte, error) {
	jss := [][]byte{}
	for ; ; jt = jt.Nested() {
		headJs, err := json.Marshal(jt.head)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		jss = append(jss, headJs)

		if !jt.Nesting() {
			break
		}
	}

	clmsJs, err := json.Marshal(jt.clms)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	jss = append(jss, clmsJs)

	return jss, nil
}
