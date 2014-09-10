package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// mongodb をバックエンドに使う。

// 非キャッシュ用。
// {
//   keyTag: "aaaa-bbbb-cccc",
//   valTag:  "https://realglobe.jp/query",
// }
type mongoKeyValueStore struct {
	*mongoDriver
	keyTag string
	valTag string
}

// keyTag と valTag は . を含まないこと。
func newMongoKeyValueStore(url, dbName, collName, keyTag, valTag string) (*mongoKeyValueStore, error) {
	base, err := newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{keyTag},
			Unique:   true,
			DropDups: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &mongoKeyValueStore{base, keyTag, valTag}, nil
}

func (reg *mongoKeyValueStore) get(key string) (value interface{}, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{reg.keyTag: key}).Select(bson.M{reg.valTag: 1, "_id": 0})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	var res map[string]interface{}
	if err := query.One(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res[reg.valTag], nil
}

func (reg *mongoKeyValueStore) put(key string, value interface{}) error {
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{reg.keyTag: key}, bson.M{reg.keyTag: key, reg.valTag: value}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoKeyValueStore) remove(key string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{reg.keyTag: key}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// キャッシュ用。
// {
//   keyTag: "aaaa-bbbb-cccc",
//   valTag:  "https://realglobe.jp/query",
//   "stamp": {
//     "date": "XXXXX",
//     "digest": "YYYYY"
//   }
// }
type mongoDatedKeyValueStore struct {
	*datedMongoDriver
	keyTag string
	valTag string
}

// keyTag と valTag は . を含まず "stamp" でないこと。
func newMongoDatedKeyValueStore(url, dbName, collName string, expiDur time.Duration, keyTag, valTag string) (*mongoDatedKeyValueStore, error) {
	reg, err := newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{keyTag},
			Unique:   true,
			DropDups: true,
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &mongoDatedKeyValueStore{newDatedMongoDriver(reg, expiDur), keyTag, valTag}, nil
}

func (reg *mongoDatedKeyValueStore) stampedGet(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{reg.keyTag: key}).Select(bson.M{reg.valTag: 1, "stamp": 1, "_id": 0})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil, nil
	}
	var res map[string]interface{}
	if err := query.One(&res); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	var date time.Time
	var dig string
	if a := res["stamp"]; a == nil {
		return nil, nil, erro.New("stamp of " + key + " is not exist.")
	} else if b, ok := a.(map[string]interface{}); !ok {
		return nil, nil, erro.New("stamp of " + key + " is invalid.")
	} else {
		if c := b["date"]; c == nil {
			return nil, nil, erro.New("date of " + key + " is not exist.")
		} else if date, ok = c.(time.Time); !ok {
			return nil, nil, erro.New("date of " + key + " is invalid.")
		}
		if d := b["digest"]; d == nil {
			return nil, nil, erro.New("digest of " + key + " is not exist.")
		} else if dig, ok = d.(string); !ok {
			return nil, nil, erro.New("digest of " + key + " is invalid.")
		}
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: date, ExpiDate: time.Now().Add(reg.expiDur), Digest: dig}
	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	return res[reg.valTag], newCaStmp, nil
}

func (reg *mongoDatedKeyValueStore) stampedPut(key string, value interface{}) (newCaStmp *Stamp, err error) {
	var digest string
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{reg.keyTag: key}).Select(bson.M{"stamp": 1, "_id": 0})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n > 0 {
		var res struct {
			Stmp *Stamp `bson:"stamp"`
		}
		if err := query.One(&res); err != nil {
			return nil, erro.Wrap(err)
		}
		n, err := strconv.Atoi(res.Stmp.Digest)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		digest = strconv.Itoa(n + 1)
	} else {
		digest = "0"
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: time.Now(), ExpiDate: time.Now().Add(reg.expiDur), Digest: digest}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{reg.keyTag: key}, bson.M{reg.keyTag: key, reg.valTag: value, "stamp": &Stamp{Date: newCaStmp.Date, Digest: newCaStmp.Digest}}); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (reg *mongoDatedKeyValueStore) remove(key string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{reg.keyTag: key}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
