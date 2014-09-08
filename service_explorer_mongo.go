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
//     "uri":  "https://realglobe.jp/api"
//   }
// }
func NewMongoServiceExplorer(url, dbName, collName string) (ServiceExplorer, error) {
	return newMongoRegistry(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"service.uri"},
			Unique:   true,
			DropDups: true,
		},
	})
}

func (reg *mongoRegistry) ServiceUuid(servUri string) (servUuid string, err error) {
	// TODO 二分探索。
	for curServUri := servUri; ; {
		query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"service.uri": curServUri})
		if n, err := query.Count(); err != nil {
			return "", erro.Wrap(err)
		} else if n > 0 {
			var res struct {
				Service struct {
					Uuid string
				}
			}
			if err := query.One(&res); err != nil {
				return "", erro.Wrap(err)
			}
			return res.Service.Uuid, nil
		}

		if serviceExplorerTreeIsRoot(curServUri) {
			break
		}

		curServUri = serviceExplorerTreeParent(curServUri)
	}
	return "", nil
}

// キャッシュ用。
// {
//   "service": {
//     "uuid": "aaaa-bbbb-cccc",
//     "uri":  "https://realglobe.jp/query"
//   },
//   "stamp": {
//     "date": "XXXXX",
//     "digest": "YYYYY"
//   }
// }
func NewMongoDatedServiceExplorer(url, dbName, collName string, expiDur time.Duration) (DatedServiceExplorer, error) {
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

func (reg *mongoBackend) StampedServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	// TODO 二分探索。
	var res struct {
		Service struct {
			Uuid string
		}
		*Stamp
	}
	for curServUri := servUri; ; {
		query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"service.uri": curServUri})
		if n, err := query.Count(); err != nil {
			return "", nil, erro.Wrap(err)
		} else if n > 0 {
			if err := query.One(&res); err != nil {
				return "", nil, erro.Wrap(err)
			}
			break
		}

		if serviceExplorerTreeIsRoot(curServUri) {
			break
		}

		curServUri = serviceExplorerTreeParent(curServUri)
	}

	if res.Stamp == nil {
		return "", nil, nil
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: res.Stamp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: res.Stamp.Digest}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return "", newCaStmp, nil
	}

	// 無効なキャッシュだった。

	return res.Service.Uuid, newCaStmp, nil
}
