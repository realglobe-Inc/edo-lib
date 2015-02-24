package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/realglobe-Inc/go-lib/erro"
	"io/ioutil"
)

func ParsePem(pemData []byte) (interface{}, error) {
	for block, rest := pem.Decode(pemData); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "PUBLIC KEY":
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return nil, erro.Wrap(err)
			}
			return key, nil
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
	return nil, erro.New("no supported key")
}

func ReadPem(path string) (interface{}, error) {
	pemData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return ParsePem(pemData)
}
