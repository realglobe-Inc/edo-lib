package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type mongoTimeLimitedKeyValueStore struct {
	*mongoDriver
	keyTag string
	valTag string
}

func NewMongoTimeLimitedKeyValueStore(url, dbName, collName, keyTag, valTag string) (TimeLimitedKeyValueStore, error) {
	base, err := newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{keyTag},
			Unique:   true,
			DropDups: true,
		},
		mgo.Index{
			Key: []string{"deadline"},
		},
	})
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &mongoTimeLimitedKeyValueStore{base, keyTag, valTag}, err
}

func (reg *mongoTimeLimitedKeyValueStore) Get(key string) (value interface{}, err error) {
	// 古いのを削除。
	reg.DB(reg.dbName).C(reg.collName).RemoveAll(bson.M{"deadline": bson.M{"$lt": time.Now()}})

	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{reg.keyTag: key})
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

func (reg *mongoTimeLimitedKeyValueStore) Put(key string, value interface{}, timLim time.Time) error {
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{reg.keyTag: key}, bson.M{reg.keyTag: key, reg.valTag: value, "deadline": timLim}); err != nil {
		return erro.Wrap(err)
	}

	// 古いのを削除。
	reg.DB(reg.dbName).C(reg.collName).RemoveAll(bson.M{"deadline": bson.M{"$lt": time.Now()}})
	return nil
}

func (reg *mongoTimeLimitedKeyValueStore) Remove(key string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{reg.keyTag: key}); err != nil {
		return erro.Wrap(err)
	}

	// 古いのを削除。
	reg.DB(reg.dbName).C(reg.collName).RemoveAll(bson.M{"deadline": bson.M{"$lt": time.Now()}})
	return nil
}
