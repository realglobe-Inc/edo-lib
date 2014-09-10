package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// mongodb をバックエンドに使う。

// 非キャッシュ用。
// {
//   "id_provider": {
//     "uuid": "aaaa-bbbb-cccc",
//     "name": "realglobe",
//     "login_uri":  "https://realglobe.jp/login"
//   }
// }
func NewMongoIdProviderLister(url, dbName, collName string) (IdProviderLister, error) {
	return newMongoDriver(url, dbName, collName, nil)
}

func (reg *mongoDriver) IdProviders() ([]*IdProvider, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"id_provider": bson.M{"$exists": true}})
	var buff []struct {
		*IdProvider `bson:"id_provider"`
	}
	if err := query.Iter().All(&buff); err != nil {
		return nil, erro.Wrap(err)
	}
	idps := []*IdProvider{}
	for _, idp := range buff {
		idps = append(idps, idp.IdProvider)
	}
	return idps, nil
}

// キャッシュ用。
// 非キャッシュ用のドキュメントに加えて、
// {
//   "stamp": {
//     "date": "XXXXX",
//     "digest": "YYYYY"
//   }
// }
func NewMongoDatedIdProviderLister(url, dbName, collName string, expiDur time.Duration) (DatedIdProviderLister, error) {
	reg, err := newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:    []string{"stamp"},
			Sparse: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newDatedMongoDriver(reg, expiDur), nil
}

func (reg *datedMongoDriver) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	var stmp *Stamp
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"stamp": bson.M{"$exists": true}})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n > 0 {
		var buff struct {
			*Stamp
		}
		if err := query.One(&buff); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		stmp = buff.Stamp
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

	query = reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"id_provider": bson.M{"$exists": true}})
	var buff []struct {
		*IdProvider `bson:"id_provider"`
	}
	if err := query.Iter().All(&buff); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	idps := []*IdProvider{}
	for _, idp := range buff {
		idps = append(idps, idp.IdProvider)
	}
	return idps, newCaStmp, nil
}
