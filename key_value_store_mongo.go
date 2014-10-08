package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

// スレッドセーフ。

const (
	mongoKeyTag   = "key"
	mongoValueTag = "value"
	mongoStampTag = "stamp"
)

type MongoMarshal func(interface{}) (interface{}, error)
type MongoUnmarshal func(interface{}) (interface{}, error)
type MongoTake func(*mgo.Query) (interface{}, *Stamp, error)

// {
//     "key": "key-no-atai",
//     "value":  value-no-atai,
//     "stamp": {
//         "date": "date-no-atai",
//         "expiration_date": "expiration-date-no-atai",
//         "digest": "digest-no-atai"
//     }
// }

type mongoKeyValueStore struct {
	base *mongoDriver
	MongoMarshal
	MongoUnmarshal
	MongoTake
}

// スレッドセーフ。
func NewMongoKeyValueStore(url, dbName, collName string, expiDur time.Duration) (KeyValueStore, error) {
	return newMongoKeyValueStore(url, dbName, collName, expiDur)
}

// スレッドセーフ。
func newMongoKeyValueStore(url, dbName, collName string, expiDur time.Duration) (*mongoKeyValueStore, error) {
	base, err := newMongoDriver(url, dbName, collName, expiDur, []mgo.Index{
		mgo.Index{
			Key:      []string{mongoKeyTag},
			Unique:   true,
			DropDups: true,
		},
		mgo.Index{
			Key: []string{mongoStampTag + ".expiration_date"},
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &mongoKeyValueStore{base: base}, nil
}

func (reg *mongoKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	query := reg.base.C().Find(bson.M{mongoKeyTag: key}).Select(bson.M{mongoValueTag: 1, mongoStampTag: 1})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil, nil
	}

	var stmp *Stamp
	if reg.MongoTake != nil {
		value, stmp, err = reg.MongoTake(query)
		if err != nil {
			return nil, nil, erro.Wrap(err)
		}
	} else {
		var res struct {
			Value interface{}
			Stamp *Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		value = res.Value
		stmp = res.Stamp
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.base.expiDur), Digest: stmp.Digest}
	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	if reg.MongoUnmarshal == nil {
		return value, newCaStmp, nil
	}

	value, err = reg.MongoUnmarshal(value)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return value, newCaStmp, nil
}

func (reg *mongoKeyValueStore) Put(key string, value interface{}) (newCaStmp *Stamp, err error) {
	var digest string
	query := reg.base.C().Find(bson.M{mongoKeyTag: key}).Select(bson.M{mongoStampTag: 1})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n > 0 {
		var res struct {
			Stamp *Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, erro.Wrap(err)
		}
		n, err := strconv.Atoi(res.Stamp.Digest)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		digest = strconv.Itoa(n + 1)
	} else {
		digest = "0"
	}

	// 対象のスタンプを取得。

	newCaStmp = &Stamp{Date: time.Now(), ExpiDate: time.Now().Add(reg.base.expiDur), Digest: digest}

	var buff interface{}
	if reg.MongoMarshal == nil {
		buff = value
	} else {
		var err error
		buff, err = reg.MongoMarshal(value)
		if err != nil {
			return nil, erro.Wrap(err)
		}
	}

	if _, err := reg.base.C().Upsert(bson.M{mongoKeyTag: key}, bson.M{mongoKeyTag: key, mongoValueTag: buff, mongoStampTag: &Stamp{Date: newCaStmp.Date, Digest: newCaStmp.Digest}}); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (reg *mongoKeyValueStore) Remove(key string) error {
	if err := reg.base.C().Remove(bson.M{mongoKeyTag: key}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
