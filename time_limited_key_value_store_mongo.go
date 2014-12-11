package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

const mongoExpiDateTag = "expiration_date"

type MongoTimeLimitedKeyValueStore interface {
	TimeLimitedKeyValueStore
	Clear() error
}

type mongoTimeLimitedKeyValueStore struct {
	base *mongoDriver

	beforeWrite Convert
	afterRead   Convert
	read        ReadDocument

	staleDur time.Duration
	expiDur  time.Duration

	date   time.Time
	digest int
}

// スレッドセーフ。
func NewMongoTimeLimitedKeyValueStore(url, dbName, collName string, beforeWrite, afterRead Convert, read ReadDocument, staleDur, expiDur time.Duration) MongoTimeLimitedKeyValueStore {
	return newMongoTimeLimitedKeyValueStore(url, dbName, collName, beforeWrite, afterRead, read, staleDur, expiDur)
}

// スレッドセーフ。
func newMongoTimeLimitedKeyValueStore(url, dbName, collName string, beforeWrite, afterRead Convert, read ReadDocument, staleDur, expiDur time.Duration) *mongoTimeLimitedKeyValueStore {
	base := newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{mongoKeyTag},
			Unique:   true,
			DropDups: true,
		},
		mgo.Index{
			Key: []string{mongoStmpTag + "." + mongoExpiDateTag},
		},
	})
	if beforeWrite == nil {
		beforeWrite = func(val interface{}) (interface{}, error) { return val, nil }
	}
	if afterRead == nil {
		afterRead = func(data interface{}) (interface{}, error) { return data, nil }
	}
	if read == nil {
		read = func(query *mgo.Query) (interface{}, *Stamp, error) {
			var res struct {
				V interface{}
				S *Stamp
			}
			if err := query.One(&res); err != nil {
				return nil, nil, erro.Wrap(err)
			}
			return res.V, res.S, nil
		}
	}
	return &mongoTimeLimitedKeyValueStore{
		base:        base,
		beforeWrite: beforeWrite,
		afterRead:   afterRead,
		read:        read,
		staleDur:    staleDur,
		expiDur:     expiDur,
		date:        time.Now(),
		digest:      0,
	}
}

func (reg *mongoTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 古いのを削除。
	coll.RemoveAll(bson.M{mongoStmpTag + "." + mongoExpiDateTag: bson.M{"$lt": time.Now()}})

	query := coll.Find(bson.M{mongoKeyTag: key}).Select(bson.M{mongoValTag: 1, mongoStmpTag: 1})
	if n, err := query.Count(); err != nil {
		reg.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil, nil
	}

	val, stmp, err := reg.read(query)
	if err != nil {
		reg.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	now := time.Now()
	newCaStmp = &Stamp{
		Date:      stmp.Date,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
		Digest:    stmp.Digest,
	}
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	val, err = reg.afterRead(val)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return val, newCaStmp, nil
}

func (reg *mongoTimeLimitedKeyValueStore) Put(key string, val interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	// 古いのを削除。
	coll.RemoveAll(bson.M{mongoStmpTag + "." + mongoExpiDateTag: bson.M{"$lt": time.Now()}})

	var digest string
	query := coll.Find(bson.M{mongoKeyTag: key}).Select(bson.M{mongoStmpTag: 1})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n > 0 {
		var res struct {
			S *Stamp
		}
		if err := query.One(&res); err != nil {
			reg.base.closeIfError()
			return nil, erro.Wrap(err)
		}
		n, err := strconv.ParseInt(res.S.Digest, 16, 64)
		if err != nil {
			return nil, erro.Wrap(err)
		}
		digest = strconv.FormatInt(n+1, 16)
	} else {
		digest = "0"
	}

	// 対象のスタンプを取得。

	now := time.Now()
	newCaStmp = &Stamp{
		Date:      now,
		StaleDate: now.Add(reg.staleDur),
		ExpiDate:  now.Add(reg.expiDur),
		Digest:    digest,
	}
	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
	}

	buff, err := reg.beforeWrite(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if _, err := coll.Upsert(bson.M{mongoKeyTag: key}, bson.M{mongoKeyTag: key, mongoValTag: buff, mongoStmpTag: &Stamp{Date: newCaStmp.Date, ExpiDate: expiDate, Digest: newCaStmp.Digest}}); err != nil {
		reg.base.closeIfError()
		return nil, erro.Wrap(err)
	}
	reg.date = now
	reg.digest++
	return newCaStmp, nil
}

func (reg *mongoTimeLimitedKeyValueStore) Remove(key string) error {
	coll, err := reg.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	// 古いのを削除。
	coll.RemoveAll(bson.M{mongoStmpTag + "." + mongoExpiDateTag: bson.M{"$lt": time.Now()}})

	if err := coll.Remove(bson.M{mongoKeyTag: key}); err != nil {
		reg.base.closeIfError()
		return erro.Wrap(err)
	}
	reg.date = time.Now()
	reg.digest++
	return nil
}

func (reg *mongoTimeLimitedKeyValueStore) Clear() error {
	coll, err := reg.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.DropCollection(); err != nil {
		reg.base.closeIfError()
		return erro.Wrap(err)
	}
	reg.date = time.Now()
	reg.digest++
	return nil
}
