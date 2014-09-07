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
//     "query_uri":  "https://realglobe.jp/query"
//   }
// }
func NewMongoIdProviderRegistry(url, dbName, collName string) (IdProviderRegistry, error) {
	return newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"id_provider.uuid"},
			Unique:   true,
			DropDups: true,
		},
	})
}

func (reg *mongoRegistry) IdProviderQueryUri(idpUuid string) (queryUri string, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"id_provider.uuid": idpUuid})
	if n, err := query.Count(); err != nil {
		return "", erro.Wrap(err)
	} else if n == 0 {
		return "", nil
	}
	var res struct {
		Id_provider struct {
			Query_uri string
		}
	}
	if err := query.One(&res); err != nil {
		return "", erro.Wrap(err)
	}
	return res.Id_provider.Query_uri, nil
}

// キャッシュ用。
// {
//   "id_provider": {
//     "uuid": "aaaa-bbbb-cccc",
//     "query_uri":  "https://realglobe.jp/query"
//   },
//   "stamp": {
//     "date": "XXXXX",
//     "digest": "YYYYY"
//   }
// }
func NewMongoDatedIdProviderRegistry(url, dbName, collName string, expiDur time.Duration) (DatedIdProviderRegistry, error) {
	reg, err := newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"id_provider.uuid"},
			Unique:   true,
			DropDups: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newMongoBackend(reg, expiDur), nil
}

func (reg *mongoBackend) StampedIdProviderQueryUri(idpUuid string, caStmp *Stamp) (queryUri string, newCaStmp *Stamp, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"id_provider.uuid": idpUuid})
	if n, err := query.Count(); err != nil {
		return "", nil, erro.Wrap(err)
	} else if n == 0 {
		return "", nil, nil
	}
	var res struct {
		Id_provider struct {
			Query_uri string
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

	return res.Id_provider.Query_uri, newCaStmp, nil
}
