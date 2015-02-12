package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"math/big"
	"strings"
)

// JSON Web Signature
type Jws interface {
	Jwt
	Verify(keys map[string]crypto.PublicKey) error
	Sign(keys map[string]crypto.PrivateKey) error
}

func NewJws() Jws {
	return newJws()
}

func ParseJws(raw string) (Jws, error) {
	jw := &jws{}
	if err := jw.parse(raw); err != nil {
		return nil, erro.Wrap(err)
	}
	return jw, nil
}

type jws struct {
	jwt

	sig []byte
}

func newJws() *jws {
	return &jws{jwt: *newJwt()}
}

func (this *jws) parse(raw string) error {
	parts := strings.Split(raw, ".")
	if len(parts) < 3 {
		return erro.New("lack of JWS parts")
	} else if len(parts) > 3 {
		return erro.New("too many JWS parts")
	} else if err := this.parseParts(parts[0], parts[1], parts[2]); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (this *jws) parseParts(headPart, clmsPart, sigPart string) error {
	if err := this.jwt.parseParts(headPart, clmsPart); err != nil {
		return erro.Wrap(err)
	}

	sig, err := base64UrlDecodeString(sigPart)
	if err != nil {
		return erro.Wrap(err)
	}

	this.sig = sig
	return nil
}

func (this *jws) SetHeader(tag string, val interface{}) {
	this.sig = nil
	this.jwt.SetHeader(tag, val)
}
func (this *jws) SetClaim(tag string, val interface{}) {
	this.sig = nil
	this.jwt.SetClaim(tag, val)
}
func (this *jws) Encode() ([]byte, error) {
	if this.sig == nil {
		return nil, erro.New("not signed")
	}

	buff, err := this.jwt.Encode()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	buff = append(buff, '.')
	buff = append(buff, base64UrlEncode(this.sig)...)
	return buff, nil
}

func (this *jws) Verify(keys map[string]crypto.PublicKey) (err error) {
	alg, _ := this.Header("alg").(string)
	if alg == "none" {
		return nil
	} else if alg == "" {
		return erro.New("no alg")
	}

	kid, _ := this.Header("kid").(string)
	key := keys[kid]
	if key == nil {
		if len(keys) == 1 {
			// 1 つだけならそれを使う。
			for _, v := range keys {
				key = v
				break
			}
		} else {
			return erro.New("no verify key")
		}
	}

	switch alg {
	case "HS256", "HS384", "HS512":
		// 共有鍵方式は未対応。
		return erro.New("not supported alg " + alg)
	case "RS256":
		return this.verifyRsaPkcs(key, crypto.SHA256)
	case "RS384":
		return this.verifyRsaPkcs(key, crypto.SHA384)
	case "RS512":
		return this.verifyRsaPkcs(key, crypto.SHA512)
	case "ES256":
		// JWA の仕様で ESxxx は鍵のサイズが決められている。
		pubKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return erro.New("not ECDSA public key")
		} else if pubKey.Params().BitSize != 256 {
			return erro.New("not P-256 EC key")
		}
		return this.verifyEcdsa(pubKey, crypto.SHA256)
	case "ES384":
		pubKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return erro.New("not ECDSA public key")
		} else if pubKey.Params().BitSize != 384 {
			return erro.New("not P-384 EC key")
		}
		return this.verifyEcdsa(pubKey, crypto.SHA384)
	case "ES512":
		pubKey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return erro.New("not ECDSA public key")
		} else if pubKey.Params().BitSize != 521 {
			return erro.New("not P-521 EC key")
		}
		return this.verifyEcdsa(pubKey, crypto.SHA512)
	case "PS256":
		return this.verifyRsaPss(key, crypto.SHA256)
	case "PS384":
		return this.verifyRsaPss(key, crypto.SHA384)
	case "PS512":
		return this.verifyRsaPss(key, crypto.SHA512)
	default:
		return erro.New("invalid alg " + alg)
	}
}

func (this *jws) verifyRsaPkcs(key crypto.PublicKey, hash crypto.Hash) error {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return erro.New("not rsa public key")
	}
	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	if err := rsa.VerifyPKCS1v15(pubKey, hash, h.Sum(nil), this.sig); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (this *jws) verifyRsaPss(key crypto.PublicKey, hash crypto.Hash) error {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return erro.New("not rsa public key")
	}
	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	if err := rsa.VerifyPSS(pubKey, hash, h.Sum(nil), this.sig, &rsa.PSSOptions{hash.Size(), hash}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (this *jws) verifyEcdsa(pubKey *ecdsa.PublicKey, hash crypto.Hash) error {
	byteSize := (pubKey.Params().BitSize + 7) / 8

	if len(this.sig) != 2*byteSize {
		return erro.New("invalid sign length ", len(this.sig))
	}

	r, s := (&big.Int{}).SetBytes(this.sig[:byteSize]), (&big.Int{}).SetBytes(this.sig[byteSize:])
	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	if !ecdsa.Verify(pubKey, h.Sum(nil), r, s) {
		return erro.New("varification failed")
	}
	return nil
}

func (this *jws) Sign(keys map[string]crypto.PrivateKey) error {
	if this.sig != nil {
		return nil
	}

	alg, _ := this.Header("alg").(string)
	if alg == "none" {
		this.sig = []byte{}
		return nil
	} else if alg == "" {
		return erro.New("no alg")
	}

	kid, _ := this.Header("kid").(string)
	key := keys[kid]
	if key == nil {
		if len(keys) == 1 {
			// 1 つだけならそれを使う。
			for k, v := range keys {
				kid = k
				key = v
				break
			}
		} else {
			return erro.New("no sign key")
		}
	}

	switch alg {
	case "HS256", "HS384", "HS512":
		// 共有鍵方式は未対応。
		return erro.New("not supported alg " + alg)
	case "RS256":
		return this.signRsaPkcs(key, crypto.SHA256)
	case "RS384":
		return this.signRsaPkcs(key, crypto.SHA384)
	case "RS512":
		return this.signRsaPkcs(key, crypto.SHA512)
	case "ES256":
		// JWA の仕様で ESxxx は鍵のサイズが決められている。
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 256 {
			return erro.New("not P-256 EC key")
		}
		return this.signEcdsa(priKey, crypto.SHA256)
	case "ES384":
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 384 {
			return erro.New("not P-384 EC key")
		}
		return this.signEcdsa(priKey, crypto.SHA384)
	case "ES512":
		priKey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return erro.New("not ECDSA private key")
		} else if priKey.Params().BitSize != 521 {
			return erro.New("not P-521 EC key")
		}
		return this.signEcdsa(priKey, crypto.SHA512)
	case "PS256":
		return this.signRsaPss(key, crypto.SHA256)
	case "PS384":
		return this.signRsaPss(key, crypto.SHA384)
	case "PS512":
		return this.signRsaPss(key, crypto.SHA512)
	default:
		return erro.New("invalid alg " + alg)
	}
}

func (this *jws) signRsaPkcs(key crypto.PrivateKey, hash crypto.Hash) error {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return erro.New("not RSA private key")
	}
	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	this.sig, err = rsa.SignPKCS1v15(rand.Reader, priKey, hash, h.Sum(nil))
	if err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (this *jws) signRsaPss(key crypto.PrivateKey, hash crypto.Hash) error {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return erro.New("not RSA private key")
	}
	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	this.sig, err = rsa.SignPSS(rand.Reader, priKey, hash, h.Sum(nil), &rsa.PSSOptions{hash.Size(), hash})
	if err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (this *jws) signEcdsa(priKey *ecdsa.PrivateKey, hash crypto.Hash) error {
	byteSize := (priKey.Params().BitSize + 7) / 8

	buff, err := this.jwt.Encode()
	if err != nil {
		return erro.Wrap(err)
	}
	h := hash.New()
	h.Write(buff)
	r, s, err := ecdsa.Sign(rand.Reader, priKey, h.Sum(nil))
	if err != nil {
		return erro.Wrap(err)
	}

	sig := make([]byte, 2*byteSize)
	rBuff := r.Bytes()
	sBuff := s.Bytes()
	copy(sig[(byteSize-len(rBuff)):byteSize], rBuff)
	copy(sig[byteSize+(byteSize-len(sBuff)):], sBuff)
	this.sig = sig
	return nil
}
