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
	"crypto/rand"
	"crypto/rsa"

	"github.com/realglobe-Inc/edo-lib/jwk"
	"github.com/realglobe-Inc/edo-lib/secrand"
	"github.com/realglobe-Inc/go-lib/erro"
)

func aCbcHsEncrypt(key []byte, keySize int, hGen crypto.Hash, plain, authData []byte) (initVec, enced, authTag []byte, err error) {
	if len(key) != keySize {
		return nil, nil, nil, erro.New("key size ", len(key), " is not ", keySize)
	}

	initVec, err = secrand.Bytes(16)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	enced, authTag, err = encryptAesCbcHmacSha2(key, hGen, plain, authData, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, enced, authTag, nil
}

func aCbcHsDecrypt(key []byte, keySize int, hGen crypto.Hash, authData, initVec, enced, authTag []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, erro.New("key size ", len(key), " is not ", keySize)
	}

	return decryptAesCbcHmacSha2(key, hGen, authData, initVec, enced, authTag)
}

func aGcmEncrypt(key []byte, keySize int, plain, authData []byte) (initVec, enced, authTag []byte, err error) {
	if len(key) != keySize {
		return nil, nil, nil, erro.New("key size ", len(key), " is not ", keySize)
	}

	initVec, err = secrand.Bytes(12)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	enced, authTag, err = encryptAesGcm(key, plain, authData, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, enced, authTag, nil
}

func aGcmDecrypt(key []byte, keySize int, authData, initVec, enced, authTag []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, erro.New("key size ", len(key), " is not ", keySize)
	}

	return decryptAesGcm(key, authData, initVec, enced, authTag)
}

// 以下、鍵交換用。

func rsa15Encrypt(key jwk.Key, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return rsa.EncryptPKCS1v15(rand.Reader, key.Public().(*rsa.PublicKey), plain)
}

func rsa15Decrypt(key jwk.Key, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return rsa.DecryptPKCS1v15(rand.Reader, key.Private().(*rsa.PrivateKey), enced)
}

func rsaOaepEncrypt(key jwk.Key, hGen crypto.Hash, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return rsa.EncryptOAEP(hGen.New(), rand.Reader, key.Public().(*rsa.PublicKey), plain, nil)
}

func rsaOaepDecrypt(key jwk.Key, hGen crypto.Hash, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return rsa.DecryptOAEP(hGen.New(), rand.Reader, key.Private().(*rsa.PrivateKey), enced, nil)
}

func aKwEncrypt(key jwk.Key, keySize int, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	} else if len(key.Common()) != keySize {
		return nil, erro.New("key size ", len(key.Common()), " is not ", keySize)
	}
	return encryptAesKw(key.Common(), plain)
}

func aKwDecrypt(key jwk.Key, keySize int, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	} else if len(key.Common()) != keySize {
		return nil, erro.New("key size ", key.Common(), " is not ", keySize)
	}
	return decryptAesKw(key.Common(), enced)
}

func dirDecrypt(key jwk.Key, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	} else if len(enced) != 0 {
		return nil, erro.New("invalid data")
	}
	return key.Common(), nil
}

func ecdhEsEncrypt(key jwk.Key, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}

func ecdhEsDecrypt(key jwk.Key, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}

func ecdhEsAKwEncrypt(key jwk.Key, keySize int, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}

func ecdhEsAKwDecrypt(key jwk.Key, keySize int, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}

func aGcmKwEncrypt(key jwk.Key, keySize int, plain []byte) (initVec, enced, authTag []byte, err error) {
	if key == nil {
		return nil, nil, nil, erro.New("no key")
	} else if len(key.Common()) != keySize {
		return nil, nil, nil, erro.New("key size ", len(key.Common()), " is not ", keySize)
	}

	initVec, err = secrand.Bytes(12)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	enced, authTag, err = encryptAesGcm(key.Common(), plain, nil, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, enced, authTag, nil
}

func aGcmKwDecrypt(key jwk.Key, keySize int, initVec, enced, authTag []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	} else if len(key.Common()) != keySize {
		return nil, erro.New("key size ", len(key.Common()), " is not ", keySize)
	}

	plain, err := decryptAesGcm(key.Common(), nil, initVec, enced, authTag)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return plain, nil
}

func pbes2HsAKwEncrypt(key jwk.Key, hGen crypto.Hash, keySize int, plain []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}

func pbes2HsAKwDecrypt(key jwk.Key, hGen crypto.Hash, keySize int, enced []byte) ([]byte, error) {
	if key == nil {
		return nil, erro.New("no key")
	}
	return nil, erro.New("not implemented")
}
