package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoTimeLimitedKeyValueStore interface {
	TimeLimitedKeyValueStore
	SetMongoTake(MongoTake)
}

type mongoTimeLimitedKeyValueStore struct {
	base *mongoKeyValueStore
}

// スレッドセーフ。
func NewMongoTimeLimitedKeyValueStore(url, dbName, collName string, expiDur time.Duration) (MongoTimeLimitedKeyValueStore, error) {
	return newMongoTimeLimitedKeyValueStore(url, dbName, collName, expiDur)
}

// スレッドセーフ。
func newMongoTimeLimitedKeyValueStore(url, dbName, collName string, expiDur time.Duration) (*mongoTimeLimitedKeyValueStore, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &mongoTimeLimitedKeyValueStore{base}, err
}

func (reg *mongoTimeLimitedKeyValueStore) SetMongoTake(take MongoTake) {
	reg.base.MongoTake = take
}

func (reg *mongoTimeLimitedKeyValueStore) Get(key string, caStmp *Stamp) (value interface{}, newCaStmp *Stamp, err error) {
	// 古いのを削除。
	reg.base.base.C().RemoveAll(bson.M{"stamp.expiration_date": bson.M{"$lt": time.Now()}})

	return reg.base.Get(key, caStmp)
}

func (reg *mongoTimeLimitedKeyValueStore) Put(key string, value interface{}, expiDate time.Time) (newCaStmp *Stamp, err error) {
	newCaStmp, err = reg.base.Put(key, value)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	// 期限を設定。
	if err := reg.base.base.C().Update(bson.M{mongoKeyTag: key}, bson.M{"$set": bson.M{"stamp.expiration_date": expiDate}}); err != nil {
		return nil, erro.Wrap(err)
	}

	if newCaStmp.ExpiDate.After(expiDate) {
		newCaStmp.ExpiDate = expiDate
	}

	// 古いのを削除。
	reg.base.base.C().RemoveAll(bson.M{"stamp.expiration_date": bson.M{"$lt": time.Now()}})
	return newCaStmp, nil
}

func (reg *mongoTimeLimitedKeyValueStore) Remove(key string) error {
	if err := reg.base.Remove(key); err != nil {
		return erro.Wrap(err)
	}

	// 古いのを削除。
	reg.base.base.C().RemoveAll(bson.M{"stamp.expiration_date": bson.M{"$lt": time.Now()}})
	return nil
}
