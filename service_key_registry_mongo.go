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
//   "service": {
//     "uuid": "aaaa-bbbb-cccc",
//     "public_key":  "XXXXX"
//   }
// }
func NewMongoServiceKeyRegistry(url, dbName, collName string) (ServiceKeyRegistry, error) {
	return newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"service.uuid"},
			Unique:   true,
			DropDups: true,
		},
	})
}

func (reg *mongoRegistry) ServiceKey(servUuid string) (key string, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"service.uuid": servUuid})
	if n, err := query.Count(); err != nil {
		return "", erro.Wrap(err)
	} else if n == 0 {
		return "", nil
	}
	var res struct {
		Service struct {
			Public_key string
		}
	}
	if err := query.One(&res); err != nil {
		return "", erro.Wrap(err)
	}
	return res.Service.Public_key, nil
}

// キャッシュ用。
// {
//   "service": {
//     "uuid": "aaaa-bbbb-cccc",
//     "public_key":  "XXXXX"
//   },
//   "stamp": {
//     "date": "YYYYY",
//     "digest": "ZZZZZ"
//   }
// }
func NewMongoDatedServiceKeyRegistry(url, dbName, collName string, expiDur time.Duration) (DatedServiceKeyRegistry, error) {
	reg, err := newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"service.uuid"},
			Unique:   true,
			DropDups: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newMongoBackend(reg, expiDur), nil
}

func (reg *mongoBackend) StampedServiceKey(servUuid string, caStmp *Stamp) (key string, newCaStmp *Stamp, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"service.uuid": servUuid})
	if n, err := query.Count(); err != nil {
		return "", nil, erro.Wrap(err)
	} else if n == 0 {
		return "", nil, nil
	}
	var res struct {
		Service struct {
			Public_key string
		}
		Stmp *Stamp `bson:"stamp"`
	}
	if err := query.One(&res); err != nil {
		return "", nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: res.Stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: res.Stmp.Digest}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return "", newCaStmp, nil
	}

	// 無効なキャッシュだった。

	return res.Service.Public_key, newCaStmp, nil
}
