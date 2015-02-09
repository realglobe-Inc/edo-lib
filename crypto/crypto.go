package crypto

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
)

func ParsePublicKey(pemData []byte) (crypto.PublicKey, error) {
	for block, rest := pem.Decode(pemData); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "PUBLIC KEY":
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, erro.Wrap(err)
			}
			return key, nil
		}
	}
	return nil, erro.New("no supported public key")
}

func ParsePrivateKey(pemData []byte) (crypto.PrivateKey, error) {
	for block, rest := pem.Decode(pemData); block != nil; block, rest = pem.Decode(rest) {
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
		}
	}
	return nil, erro.New("no supported private key")
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
