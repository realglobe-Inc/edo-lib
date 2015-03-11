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
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"github.com/realglobe-Inc/go-lib/erro"
	"math/big"
)

var jwsAlgs = map[string]bool{
	"none":  true,
	"HS256": true,
	"HS384": true,
	"HS512": true,
	"RS256": true,
	"RS384": true,
	"RS512": true,
	"ES256": true,
	"ES384": true,
	"ES512": true,
	"PS256": true,
	"PS384": true,
	"PS512": true,
}

func isJwsAlgorithm(alg string) bool {
	return jwsAlgs[alg]
}

func HashFunction(alg string) (crypto.Hash, error) {
	switch alg {
	case "none":
		return 0, nil
	case "HS256", "RS256", "ES256", "PS256":
		return crypto.SHA256, nil
	case "HS384", "RS384", "ES384", "PS384":
		return crypto.SHA384, nil
	case "HS512", "RS512", "ES512", "PS512":
		return crypto.SHA512, nil
	default:
		return 0, erro.New("alg " + alg + " is unsupported")
	}
}

func noneVerify(sig []byte) error {
	if len(sig) != 0 {
		return erro.New("verification failed")
	}
	return nil
}

func hsSign(key interface{}, hGen crypto.Hash, ds ...[]byte) ([]byte, error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	}

	h := hmac.New(hGen.New, comKey)
	for _, d := range ds {
		h.Write(d)
	}
	return h.Sum(nil), nil
}

func hsVerify(key interface{}, hGen crypto.Hash, sig []byte, ds ...[]byte) error {
	comKey, ok := key.([]byte)
	if !ok {
		return erro.New("not common key")
	}

	h := hmac.New(hGen.New, comKey)
	for _, d := range ds {
		h.Write(d)
	}
	if !hmac.Equal(h.Sum(nil), sig) {
		return erro.New("verification failed")
	}
	return nil
}

func rsSign(key interface{}, hGen crypto.Hash, ds ...[]byte) ([]byte, error) {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, erro.New("not RSA private key")
	}

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, priKey, hGen, h.Sum(nil))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sig, nil
}

func rsVerify(key interface{}, hGen crypto.Hash, sig []byte, ds ...[]byte) error {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return erro.New("not RSA public key")
	}

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}
	return rsa.VerifyPKCS1v15(pubKey, hGen, h.Sum(nil), sig)
}

func esSign(priKey *ecdsa.PrivateKey, hGen crypto.Hash, ds ...[]byte) ([]byte, error) {

	byteSize := (priKey.Params().BitSize + 7) / 8

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}

	r, s, err := ecdsa.Sign(rand.Reader, priKey, h.Sum(nil))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	sig := make([]byte, 2*byteSize)
	rBuff := r.Bytes()
	sBuff := s.Bytes()
	copy(sig[(byteSize-len(rBuff)):byteSize], rBuff)
	copy(sig[byteSize+(byteSize-len(sBuff)):], sBuff)
	return sig, nil
}

func esVerify(pubKey *ecdsa.PublicKey, hGen crypto.Hash, sig []byte, ds ...[]byte) error {
	byteSize := (pubKey.Params().BitSize + 7) / 8
	if len(sig) != 2*byteSize {
		return erro.New("verification failed")
	}
	r, s := (&big.Int{}).SetBytes(sig[:byteSize]), (&big.Int{}).SetBytes(sig[byteSize:])

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}

	if !ecdsa.Verify(pubKey, h.Sum(nil), r, s) {
		return erro.New("verification failed")
	}
	return nil
}

func psSign(key interface{}, hGen crypto.Hash, ds ...[]byte) ([]byte, error) {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, erro.New("not RSA private key")
	}

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}
	sig, err := rsa.SignPSS(rand.Reader, priKey, hGen, h.Sum(nil), &rsa.PSSOptions{hGen.Size(), hGen})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sig, nil
}

func psVerify(key crypto.PublicKey, hGen crypto.Hash, sig []byte, ds ...[]byte) error {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return erro.New("not RSA public key")
	}

	h := hGen.New()
	for _, d := range ds {
		h.Write(d)
	}
	return rsa.VerifyPSS(pubKey, hGen, h.Sum(nil), sig, &rsa.PSSOptions{hGen.Size(), hGen})
}
