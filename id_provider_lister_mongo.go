package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// mongodb をバックエンドに使う。

// 非キャッシュ用。
func NewMongoIdProviderLister(url, dbName, collName string) (IdProviderLister, error) {
	return newMongoRegistry(url, dbName, collName, nil)
}

func (reg *mongoRegistry) IdProviders() ([]*IdProvider, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(nil)
	var idps []*IdProvider
	if err := query.Iter().All(&idps); err != nil {
		return nil, erro.Wrap(err)
	}
	return idps, nil
}

// キャッシュ用。
func NewMongoDatedIdProviderLister(url, dbName, collName string, expiDur time.Duration) (DatedIdProviderLister, error) {
	reg, err := newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:    []string{"date"},
			Sparse: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newMongoBackend(reg, expiDur), nil
}

func (reg *mongoBackend) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	var stmp *Stamp
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"date": bson.M{"$exists": true}})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n != 0 {
		stmp = &Stamp{}
		if err := query.One(stmp); err != nil {
			return nil, nil, erro.Wrap(err)
		}
	}

	// 対象のスタンプを取得。

	var newCaStmp *Stamp
	if stmp != nil {
		newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}
	} else {
		newCaStmp = &Stamp{ExpiDate: time.Now().Add(reg.expiDur)}
	}

	if caStmp != nil && stmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	query = reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"uuid": bson.M{"$exists": true}})
	var idps []*IdProvider
	if err := query.Iter().All(&idps); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return idps, newCaStmp, nil
}
