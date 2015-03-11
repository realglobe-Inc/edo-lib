// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"time"
)

// 値自体がキーとタイムスタンプを含む必要がある。

type Convert func(interface{}) (interface{}, error)
type ReadDocument func(*mgo.Query) (interface{}, error)
type GetStamp func(interface{}) *Stamp

type MongoKeyValueStore interface {
	KeyValueStore

	// コレクションを消す。主にテスト用。
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

func (drv *mongoKeyValueStore) getStamp(val interface{}) *Stamp {
	now := time.Now()
	stmp := drv.getStmp(val)
	stmp.StaleDate = now.Add(drv.staleDur)
	stmp.ExpiDate = now.Add(drv.expiDur)
	return stmp
}

func (drv *mongoKeyValueStore) Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	coll, err := drv.base.collection()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	query := coll.Find(bson.M{drv.keyTag: key})
	val, err = drv.read(query)
	if err != nil {
		if erro.Unwrap(err) == mgo.ErrNotFound {
			return nil, nil, nil
		}
		drv.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	}
	val, err = drv.afterRd(val)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	newCaStmp = drv.getStmp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return val, newCaStmp, nil
}

func (drv *mongoKeyValueStore) Put(key string, val interface{}) (newCaStmp *Stamp, err error) {
	coll, err := drv.base.collection()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStmp(val)

	buff, err := drv.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if _, err := coll.Upsert(bson.M{drv.keyTag: key}, buff); err != nil {
		drv.base.closeIfError()
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (drv *mongoKeyValueStore) Remove(key string) error {
	coll, err := drv.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.Remove(bson.M{drv.keyTag: key}); err != nil {
		drv.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}

func (drv *mongoKeyValueStore) Close() error {
	return drv.base.close()
}

func (drv *mongoKeyValueStore) Clear() error {
	coll, err := drv.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.DropCollection(); err != nil {
		drv.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}

type NKeyValueStore interface {
	NGet(tagKeys bson.M, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error)
	NPut(tagKeys bson.M, val interface{}) (*Stamp, error)
	NRemove(tagKeys bson.M) error
	io.Closer
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

func (drv *mongoKeyValueStore) NGet(tagKeys bson.M, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	coll, err := drv.base.collection()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	query := coll.Find(tagKeys)
	val, err = drv.read(query)
	if err != nil {
		if erro.Unwrap(err) == mgo.ErrNotFound {
			return nil, nil, nil
		}
		drv.base.closeIfError()
		return nil, nil, erro.Wrap(err)
	}
	val, err = drv.afterRd(val)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 対象のスタンプを取得。

	newCaStmp = drv.getStmp(val)
	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	return val, newCaStmp, nil
}

func (drv *mongoKeyValueStore) NPut(tagKeys bson.M, val interface{}) (newCaStmp *Stamp, err error) {
	coll, err := drv.base.collection()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	newCaStmp = drv.getStmp(val)

	buff, err := drv.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if _, err := coll.Upsert(tagKeys, buff); err != nil {
		drv.base.closeIfError()
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (drv *mongoKeyValueStore) NRemove(tagKeys bson.M) error {
	coll, err := drv.base.collection()
	if err != nil {
		return erro.Wrap(err)
	}

	if err := coll.Remove(tagKeys); err != nil {
		drv.base.closeIfError()
		return erro.Wrap(err)
	}
	return nil
}
