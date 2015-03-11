package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/realglobe-Inc/edo-lib/secrand"
	"github.com/realglobe-Inc/go-lib/erro"
)

var jweAlgs = map[string]bool{
	"RSA1_5":             true,
	"RSA-OAEP":           true,
	"RSA-OAEP-256":       true,
	"A128KW":             true,
	"A192KW":             true,
	"A256KW":             true,
	"dir":                true,
	"ECDH-ES":            true,
	"ECDH-ES+A128KW":     true,
	"ECDH-ES+A192KW":     true,
	"ECDH-ES+A256KW":     true,
	"A128GCMKW":          true,
	"A192GCMKW":          true,
	"A256GCMKW":          true,
	"PBES2-HS256+A128KW": true,
	"PBES2-HS384+A192KW": true,
	"PBES2-HS512+A256KW": true,
}

var jweEncs = map[string]bool{
	"A128CBC-HS256": true,
	"A192CBC-HS384": true,
	"A256CBC-HS512": true,
	"A128GCM":       true,
	"A192GCM":       true,
	"A256GCM":       true,
}

var jweZips = map[string]bool{
	"DEF": true,
}

func isJweAlgorithm(alg string) bool {
	return jweAlgs[alg]
}

func aCbcHsEncrypt(key []byte, keySize int, hGen crypto.Hash, plain, authData []byte) (initVec, encrypted, authTag []byte, err error) {
	if len(key) != keySize {
		return nil, nil, nil, erro.New("key size is not ", keySize)
	}

	initVec, err = secrand.Bytes(16)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	encrypted, authTag, err = encryptAesCbcHmacSha2(key, hGen, plain, authData, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, encrypted, authTag, nil
}

func aCbcHsDecrypt(key []byte, keySize int, hGen crypto.Hash, authData, initVec, encrypted, authTag []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, erro.New("key size is not ", keySize)
	}

	return decryptAesCbcHmacSha2(key, hGen, authData, initVec, encrypted, authTag)
}

func aGcmEncrypt(key []byte, keySize int, plain, authData []byte) (initVec, encrypted, authTag []byte, err error) {
	if len(key) != keySize {
		return nil, nil, nil, erro.New("key size is not ", keySize)
	}

	initVec, err = secrand.Bytes(12)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	encrypted, authTag, err = encryptAesGcm(key, plain, authData, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, encrypted, authTag, nil
}

func aGcmDecrypt(key []byte, keySize int, authData, initVec, encrypted, authTag []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, erro.New("key size is not ", keySize)
	}

	return decryptAesGcm(key, authData, initVec, encrypted, authTag)
}

func rsa15Encrypt(key interface{}, plain []byte) ([]byte, error) {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, erro.New("not RSA public key")
	}

	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, plain)
}

func rsa15Decrypt(key interface{}, encrypted []byte) ([]byte, error) {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, erro.New("not RSA private key")
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priKey, encrypted)
}

func rsaOaepEncrypt(key interface{}, hGen crypto.Hash, plain []byte) ([]byte, error) {
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, erro.New("not RSA public key")
	}

	return rsa.EncryptOAEP(hGen.New(), rand.Reader, pubKey, plain, nil)
}

func rsaOaepDecrypt(key interface{}, hGen crypto.Hash, encrypted []byte) ([]byte, error) {
	priKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, erro.New("not RSA private key")
	}

	return rsa.DecryptOAEP(hGen.New(), rand.Reader, priKey, encrypted, nil)
}

func aKwEncrypt(key interface{}, keySize int, plain []byte) ([]byte, error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	} else if len(comKey) != keySize {
		return nil, erro.New("common key size is not ", keySize)
	}

	return encryptAesKw(comKey, plain)
}

func aKwDecrypt(key interface{}, keySize int, encrypted []byte) ([]byte, error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	} else if len(comKey) != keySize {
		return nil, erro.New("common key size is not ", keySize)
	}

	return decryptAesKw(comKey, encrypted)
}

func dirEncrypt(key interface{}) ([]byte, error) {
	_, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	}
	return []byte{}, nil
}

func dirDecrypt(key interface{}, encrypted []byte) ([]byte, error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	} else if len(encrypted) != 0 {
		return nil, erro.New("not empty data")
	}
	return comKey, nil
}

func ecdhEsEncrypt(key interface{}, plain []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}

func ecdhEsDecrypt(key interface{}, encrypted []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}

func ecdhEsAKwEncrypt(key interface{}, keySize int, plain []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}

func ecdhEsAKwDecrypt(key interface{}, keySize int, encrypted []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}

func aGcmKwEncrypt(key interface{}, keySize int, plain []byte) (initVec, encrypted, authTag []byte, err error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, nil, nil, erro.New("not common key")
	} else if len(comKey) != keySize {
		return nil, nil, nil, erro.New("common key size is not ", keySize)
	}

	initVec, err = secrand.Bytes(12)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}

	encrypted, authTag, err = encryptAesGcm(comKey, plain, nil, initVec)
	if err != nil {
		return nil, nil, nil, erro.Wrap(err)
	}
	return initVec, encrypted, authTag, nil
}

func aGcmKwDecrypt(key interface{}, keySize int, initVec, encrypted, authTag []byte) ([]byte, error) {
	comKey, ok := key.([]byte)
	if !ok {
		return nil, erro.New("not common key")
	} else if len(comKey) != keySize {
		return nil, erro.New("common key size is not ", keySize)
	}

	plain, err := decryptAesGcm(comKey, nil, initVec, encrypted, authTag)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return plain, nil
}

func pbes2HsAKwEncrypt(key interface{}, hGen crypto.Hash, keySize int, plain []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}

func pbes2HsAKwDecrypt(key interface{}, hGen crypto.Hash, keySize int, encrypted []byte) ([]byte, error) {
	return nil, erro.New("not implemented")
}
