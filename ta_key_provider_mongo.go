package driver

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

func publicKeyToPem(pubKey interface{}) (pemStr interface{}, err error) {
	block := &pem.Block{
		Type: "PUBLIC KEY",
	}
	block.Bytes, err = x509.MarshalPKIXPublicKey(pubKey.(*rsa.PublicKey))
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return pem.EncodeToMemory(block), nil
}

// PEM 形式の文字列から rsa.PublicKey をつくる。
func pemToPublicKey(pemStr interface{}) (pubKey interface{}, err error) {
	return util.ParseRsaPublicKey(string(pemStr.([]byte)))
}

// スレッドセーフ。
func NewMongoTaKeyProvider(url, dbName, collName string, expiDur time.Duration) (TaKeyProvider, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.MongoMarshal = publicKeyToPem
	base.MongoUnmarshal = pemToPublicKey
	// デコード後をキャッシュ。
	// TODO キャッシュの並列化。
	return newTaKeyProvider(newCachingKeyValueStore(base)), nil
}
