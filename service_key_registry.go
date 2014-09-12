package driver

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/realglobe-Inc/go-lib-rg/erro"
)

// サービスの公開鍵を返す。
type ServiceKeyRegistry interface {
	ServiceKey(servUuid string) (servKey *rsa.PublicKey, err error)
}

// サービスの公開鍵を返す。キャッシュ用。
type DatedServiceKeyRegistry interface {
	StampedServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error)
}

// KVS の中に pem 形式で保存してあるとする。

// 非キャッシュ用。
type serviceKeyRegistry struct {
	keyValueStore
}

func newServiceKeyRegistry(base keyValueStore) *serviceKeyRegistry {
	return &serviceKeyRegistry{base}
}

func parseKey(pemStr string) (*rsa.PublicKey, error) {
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
		return nil, erro.New("not public key block.")
	}
}

func (reg *serviceKeyRegistry) ServiceKey(servUuid string) (servKey *rsa.PublicKey, err error) {
	val, err := reg.get(servUuid)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if val == nil || val == "" {
		return nil, nil
	}
	return parseKey(val.(string))
}

// キャッシュ用。
type datedServiceKeyRegistry struct {
	datedKeyValueStore
}

func newDatedServiceKeyRegistry(base datedKeyValueStore) *datedServiceKeyRegistry {
	return &datedServiceKeyRegistry{base}
}

func (reg *datedServiceKeyRegistry) StampedServiceKey(servUuid string, caStmp *Stamp) (servKey *rsa.PublicKey, newCaStmp *Stamp, err error) {
	val, newCaStmp, err := reg.stampedGet(servUuid, caStmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if newCaStmp == nil {
		return nil, nil, nil
	} else if val == nil || val == "" {
		return nil, newCaStmp, nil
	}

	servKey, err = parseKey(val.(string))
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return servKey, newCaStmp, nil
}
