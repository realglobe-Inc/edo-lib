package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"reflect"
	"time"
)

// mondodb を使うドライバー。

// 非キャッシュ用。
type mongoDriver struct {
	dbName   string
	collName string
	*mgo.Session
}

func newMongoDriver(url, dbName, collName string, indices []mgo.Index) (*mongoDriver, error) {
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	curIndices, err := sess.DB(dbName).C(collName).Indexes()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	// 既存の要らない索引を消す。
	for _, curIdx := range curIndices {
		if len(curIdx.Key) == 1 && curIdx.Key[0] == "_id" {
			continue
		}

		ok := false
		for _, idx := range indices {
			if reflect.DeepEqual(curIdx, idx) {
				ok = true
				break
			}
		}
		if ok {
			continue
		}

		// 要らない。
		if err := sess.DB(dbName).C(collName).DropIndex(curIdx.Key...); err != nil {
			return nil, erro.Wrap(err)
		}
	}

	curIndices, err = sess.DB(dbName).C(collName).Indexes()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, idx := range indices {
		ok := true
		for _, curIdx := range curIndices {
			if reflect.DeepEqual(idx, curIdx) {
				ok = false
				break
			}
		}
		if !ok {
			// もうある。
			continue
		}

		if err := sess.DB(dbName).C(collName).EnsureIndex(idx); err != nil {
			return nil, erro.Wrap(err)
		}
	}

	return &mongoDriver{dbName, collName, sess}, nil
}

// キャッシュ用。
type datedMongoDriver struct {
	*mongoDriver
	expiDur time.Duration
}

func newDatedMongoDriver(base *mongoDriver, expiDur time.Duration) *datedMongoDriver {
	return &datedMongoDriver{base, expiDur}
}
