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
func NewMongoKeyValueStore(sess *mgo.Session, dbName, collName, keyTag string, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) MongoKeyValueStore {
	return newMongoKeyValueStore(sess, dbName, collName, keyTag, beforeWr, afterRd, read, getStmp, staleDur, expiDur)
}

// スレッドセーフ。
func newMongoKeyValueStore(sess *mgo.Session, dbName, collName, keyTag string, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) *mongoKeyValueStore {
	base := newMongoDriver(sess, dbName, collName)
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
	val, err = func() (interface{}, error) {
		sess, coll := drv.base.collection()
		defer sess.Close() // パニックでも解放するように defer、使ったらすぐ解放するように無名関数。

		raw, err := drv.read(coll.Find(bson.M{drv.keyTag: key}))
		if err != nil {
			if erro.Unwrap(err) == mgo.ErrNotFound {
				return nil, nil
			}
			return nil, erro.Wrap(err)
		}
		return drv.afterRd(raw)
	}()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if val == nil {
		return nil, nil, nil
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
	newCaStmp = drv.getStmp(val)

	buff, err := drv.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	sess, coll := drv.base.collection()
	defer sess.Close()

	if _, err := coll.Upsert(bson.M{drv.keyTag: key}, buff); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (drv *mongoKeyValueStore) Remove(key string) error {
	sess, coll := drv.base.collection()
	defer sess.Close()

	if err := coll.Remove(bson.M{drv.keyTag: key}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (drv *mongoKeyValueStore) Close() error {
	drv.base = nil
	return nil
}

func (drv *mongoKeyValueStore) Clear() error {
	sess, coll := drv.base.collection()
	defer sess.Close()

	if err := coll.DropCollection(); err != nil {
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
func NewMongoNKeyValueStore(sess *mgo.Session, dbName, collName string, tags []string, beforeWr, afterRd Convert, read ReadDocument, getStmp GetStamp, staleDur, expiDur time.Duration) MongoNKeyValueStore {
	return newMongoKeyValueStore(sess, dbName, collName, "", beforeWr, afterRd, read, getStmp, staleDur, expiDur)
}

func (drv *mongoKeyValueStore) NGet(tagKeys bson.M, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error) {
	val, err = func() (interface{}, error) {
		sess, coll := drv.base.collection()
		defer sess.Close() // パニックでも解放するように defer、使ったらすぐ解放するように無名関数。

		raw, err := drv.read(coll.Find(tagKeys))
		if err != nil {
			if erro.Unwrap(err) == mgo.ErrNotFound {
				return nil, nil
			}
			return nil, erro.Wrap(err)
		}

		return drv.afterRd(raw)
	}()
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if val == nil {
		return nil, nil, nil
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
	newCaStmp = drv.getStmp(val)

	buff, err := drv.beforeWr(val)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	sess, coll := drv.base.collection()
	defer sess.Close()

	if _, err := coll.Upsert(tagKeys, buff); err != nil {
		return nil, erro.Wrap(err)
	}
	return newCaStmp, nil
}

func (drv *mongoKeyValueStore) NRemove(tagKeys bson.M) error {
	sess, coll := drv.base.collection()
	defer sess.Close()

	if err := coll.Remove(tagKeys); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
