package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 値自体がキーとタイムスタンプを含む必要がある。

type Convert func(interface{}) (interface{}, error)
type ReadDocument func(*mgo.Query) (interface{}, error)
type GetStamp func(interface{}) *Stamp

type MongoKeyValueStore interface {
	KeyValueStore
	Clear() error
}

type mongoKeyValueStore struct {
	base *mongoDriver

	keyTag   string
	beforeWr Convert
	afterRd  Convert
	read     ReadDocument
	getStmp  GetStamp

	staleDur time.Duration
	expiDur  time.Duration
}

// スレッドセーフ。
func NewMongoKeyValueStore(url, dbName, collName, keyTag string, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) MongoKeyValueStore {
	return newMongoKeyValueStore(url, dbName, collName, keyTag, []mgo.Index{
		mgo.Index{
			Key:      []string{keyTag},
			Unique:   true,
			DropDups: true,
		},
	}, beforeWr, afterRd, read, getStmp, staleDur, expiDur)
}

// スレッドセーフ。
func newMongoKeyValueStore(url, dbName, collName, keyTag string, indices []mgo.Index, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) *mongoKeyValueStore {
	base := newMongoDriver(url, dbName, collName, indices)
	if beforeWr == nil {
		beforeWr = func(val interface{}) (interface{}, error) { return val, nil }
	}
	if afterRd == nil {
		afterRd = func(data interface{}) (interface{}, error) { return data, nil }
	}
	if read == nil {
		read = func(query *mgo.Query) (interface{}, error) {
			var res map[string]interface{}
			if err := query.One(&res); err != nil {
				return nil, erro.Wrap(err)
			}
			return res, nil
		}
	}
	if getStmp == nil {
		getStmp = func(val interface{}) *Stamp {
			m, _ := val.(map[string]interface{})
			date, _ := m["date"].(time.Time)
			dig, _ := m["digest"].(string)
			return &Stamp{Date: date, Digest: dig}
		}
	}
	return &mongoKeyValueStore{
		base:     base,
		keyTag:   keyTag,
		beforeWr: beforeWr,
		afterRd:  afterRd,
		read:     read,
		getStmp:  getStmp,
		staleDur: staleDur,
		expiDur:  expiDur,
	}
}

func (this *mongoKeyValueStore) getStamp(val interface{}) *Stamp {
	now := time.Now()
	stmp := this.getStmp(val)
	stmp.StaleDate = now.Add(this.staleDur)
	stmp.ExpiDate = now.Add(this.expiDur)
	return stmp
}

func (reg *mongoKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	query := coll.Find(bson.M{reg.keyTag: key})
	val, err = reg.read(query)
	if err != nil {
		if erro.Unwrap(err) == mgo.ErrNotFound {
			return nil, nil, nil
		}
		reg.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	}
	val, err = reg.afterRd(val)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	newCaStmp = reg.getStmp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return val, newCaStmp, nil
}

func (reg *mongoKeyValueStore) Put(key string, val interface{}) (newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp = reg.getStmp(val)

	buff, err := reg.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if _, err := coll.Upsert(bson.M{reg.keyTag: key}, buff); err != nil {
		reg.base.closeIfError()
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (reg *mongoKeyValueStore) Remove(key string) error {
	coll, err := reg.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.Remove(bson.M{reg.keyTag: key}); err != nil {
		reg.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoKeyValueStore) Clear() error {
	coll, err := reg.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.DropCollection(); err != nil {
		reg.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}

type NKeyValueStore interface {
	NGet(tagKeys bson.M, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error)
	NPut(tagKeys bson.M, val interface{}) (*Stamp, error)
	NRemove(tagKeys bson.M) error
}

type MongoNKeyValueStore interface {
	NKeyValueStore
	Clear() error
}

// スレッドセーフ。
func NewMongoNKeyValueStore(url, dbName, collName string, tags []string, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) MongoNKeyValueStore {
	return newMongoKeyValueStore(url, dbName, collName, "", []mgo.Index{
		mgo.Index{
			Key:      tags,
			Unique:   true,
			DropDups: true,
		},
	}, beforeWr, afterRd, read, getStmp, staleDur, expiDur)
}

func (reg *mongoKeyValueStore) NGet(tagKeys bson.M, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	query := coll.Find(tagKeys)
	val, err = reg.read(query)
	if err != nil {
		if erro.Unwrap(err) == mgo.ErrNotFound {
			return nil, nil, nil
		}
		reg.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	}
	val, err = reg.afterRd(val)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	newCaStmp = reg.getStmp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return val, newCaStmp, nil
}

func (reg *mongoKeyValueStore) NPut(tagKeys bson.M, val interface{}) (newCaStmp *Stamp, err error) {
	coll, err := reg.base.collection()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp = reg.getStmp(val)

	buff, err := reg.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if _, err := coll.Upsert(tagKeys, buff); err != nil {
		reg.base.closeIfError()
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (reg *mongoKeyValueStore) NRemove(tagKeys bson.M) error {
	coll, err := reg.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.Remove(tagKeys); err != nil {
		reg.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}
