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
	hash "github.com/realglobe-Inc/edo-lib/hash"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"github.com/realglobe-Inc/go-lib/erro"
	"math/big"
)

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
		return 0, erro.New("unsupported algorithm " + alg)
	}
}

func noneVerify(sig []byte) error {
	if len(sig) != 0 {
		return erro.New("verification failed")
	}
	return nil
}

func hsSign(key jwk.Key, hGen crypto.Hash, data []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return hash.Hashing(hmac.New(hGen.New, key.Common()), data), nil
}

func hsVerify(key jwk.Key, hGen crypto.Hash, sig []byte, data []byte) error {
	if key == nil {
		return erro.New("no key")
	}
	h := hash.Hashing(hmac.New(hGen.New, key.Common()), data)
	if !hmac.Equal(h, sig) {
		return erro.New("verification failed")
	}
	return nil
}

func rsSign(key jwk.Key, hGen crypto.Hash, data []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}

	sig, err := rsa.SignPKCS1v15(rand.Reader, key.Private().(*rsa.PrivateKey), hGen, hash.Hashing(hGen.New(), data))
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sig, nil
}

func rsVerify(key jwk.Key, hGen crypto.Hash, sig []byte, data []byte) error {
	if key == nil {
		return erro.New("no key")
	}

	return rsa.VerifyPKCS1v15(key.Public().(*rsa.PublicKey), hGen, hash.Hashing(hGen.New(), data), sig)
}

func esSign(key jwk.Key, hGen crypto.Hash, data []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}

	pri := key.Private().(*ecdsa.PrivateKey)
	byteSize := (pri.Params().BitSize + 7) / 8

	r, s, err := ecdsa.Sign(rand.Reader, pri, hash.Hashing(hGen.New(), data))
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

func esVerify(key jwk.Key, hGen crypto.Hash, sig []byte, data []byte) error {
	if key == nil {
		return erro.New("no key")
	}

	pub := key.Public().(*ecdsa.PublicKey)
	byteSize := (pub.Params().BitSize + 7) / 8
	if len(sig) != 2*byteSize {
		return erro.New("verification failed")
	}
	r, s := (&big.Int{}).SetBytes(sig[:byteSize]), (&big.Int{}).SetBytes(sig[byteSize:])

	if !ecdsa.Verify(pub, hash.Hashing(hGen.New(), data), r, s) {
		return erro.New("verification failed")
	}
	return nil
}

func psSign(key jwk.Key, hGen crypto.Hash, data []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}

	sig, err := rsa.SignPSS(rand.Reader, key.Private().(*rsa.PrivateKey), hGen, hash.Hashing(hGen.New(), data), &rsa.PSSOptions{hGen.Size(), hGen})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sig, nil
}

func psVerify(key jwk.Key, hGen crypto.Hash, sig []byte, data []byte) error {
	if key == nil {
		return erro.New("no key")
	}

	return rsa.VerifyPSS(key.Public().(*rsa.PublicKey), hGen, hash.Hashing(hGen.New(), data), sig, &rsa.PSSOptions{hGen.Size(), hGen})
}
