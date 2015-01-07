package util

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"strings"
)

func ParseHashFunction(hashStr string) (crypto.Hash, error) {
	switch strings.ToUpper(hashStr) {
	case "MD4": // code.google.com/p/go.crypto/md4 が要る。
		return crypto.MD4, nil
	case "MD5":
		return crypto.MD5, nil
	case "SHA1":
		return crypto.SHA1, nil
	case "SHA224":
		return crypto.SHA224, nil
	case "SHA256":
		return crypto.SHA256, nil
	case "SHA384":
		return crypto.SHA384, nil
	case "SHA512":
		return crypto.SHA512, nil
	case "MD5SHA1": // 公式に no implementation。
		return crypto.MD5SHA1, nil
	case "RIPEMD160": // code.google.com/p/go.crypto/ripemd160 が要る。
		return crypto.RIPEMD160, nil
	default:
		return 0, erro.New("not supported hash function " + hashStr + ".")
	}
}

func HashFunctionString(hash crypto.Hash) string {
	switch hash {
	case crypto.MD4:
		return "MD4"
	case crypto.MD5:
		return "MD5"
	case crypto.SHA1:
		return "SHA1"
	case crypto.SHA224:
		return "SHA224"
	case crypto.SHA256:
		return "SHA256"
	case crypto.SHA384:
		return "SHA384"
	case crypto.SHA512:
		return "SHA512"
	case crypto.MD5SHA1:
		return "MD5SHA1"
	case crypto.RIPEMD160:
		return "RIPEMD160"
	default:
		return "unknown"
	}
}

func ParseRsaPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, erro.New("no public key.")
	}

	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		switch k := key.(type) {
		case *rsa.PublicKey:
			return k, nil
		default:
			return nil, erro.New("not rsa public key.")
		}
	default:
		return nil, erro.New("invalid public key type " + block.Type + ".")
	}
}

func ParseRsaPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, erro.New("no private key.")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return priKey, nil
	default:
		return nil, erro.New("invalid private key type " + block.Type + ".")
	}
}

func ParsePublicKey(pemData []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, erro.New("no public key.")
	}

	switch block.Type {
	case "PUBLIC KEY":
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return key, nil
	default:
		return nil, erro.New("invalid public key type " + block.Type + ".")
	}
}

func ParsePrivateKey(pemData []byte) (crypto.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, erro.New("no private key.")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return key, nil
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		return key, nil
	default:
		return nil, erro.New("invalid private key type " + block.Type + ".")
	}
}

func ReadPublicKey(path string) (crypto.PublicKey, error) {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return ParsePublicKey(pemData)
}

func ReadPrivateKey(path string) (crypto.PrivateKey, error) {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return ParsePrivateKey(pemData)
}
