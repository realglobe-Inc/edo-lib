package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strconv"
	"time"
)

// JavaScript.
func NewMongoJsBackendRegistry(url, dbName, collName string) (JsBackendRegistry, error) {
	reg, err := NewMongoJsRegistry(url, dbName, collName)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return reg.(*mongoRegistry), nil
}

func (reg *mongoRegistry) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"path": path.Join(dir, objName)})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil, nil
	}
	var res mongoObject
	if err := query.One(&res); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	stmp := &Stamp{Date: res.Date, Digest: strconv.Itoa(res.Digest)}

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: time.Now(), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	return &Object{res.Service, res.Library, res.Include, res.Code}, newCaStmp, nil
}

// ID プロバイダ。
func NewMongoIdProviderBackend(url, dbName, collName string) (IdProviderBackend, error) {
	reg, err := newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:    []string{"date"},
			Sparse: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return reg, nil
}

func (reg *mongoRegistry) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
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

	newCaStmp := &Stamp{Date: time.Now()}
	if stmp != nil {
		newCaStmp.Digest = stmp.Digest
	}

	if caStmp != nil && stmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	query = reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"id_provider_uuid": bson.M{"$exists": true}})
	var idps []*IdProvider
	if err := query.Iter().All(&idps); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return idps, newCaStmp, nil
}
