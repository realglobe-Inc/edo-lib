package driver

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"gopkg.in/mgo.v2"
	"reflect"
	"sync"
)

// 接続が切れても次は接続し直す。
// スレッドセーフ。
type mongoDriver struct {
	url     string
	db      string
	coll    string
	indices []mgo.Index

	sessLock sync.Mutex
	sess     *mgo.Session
}

func newMongoDriver(url, db, coll string, indices []mgo.Index) *mongoDriver {
	return &mongoDriver{
		url:     url,
		db:      db,
		coll:    coll,
		indices: indices,
	}
}

func (drv *mongoDriver) collection() (*mgo.Collection, error) {
	sess, err := drv.session()
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sess.DB(drv.db).C(drv.coll), nil
}

func (drv *mongoDriver) session() (*mgo.Session, error) {
	drv.sessLock.Lock()
	defer drv.sessLock.Unlock()

	if drv.sess != nil {
		return drv.sess, nil
	}

	sess, err := mgo.Dial(drv.url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	curIndices, err := sess.DB(drv.db).C(drv.coll).Indexes()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, idx := range drv.indices {
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

		if err := sess.DB(drv.db).C(drv.coll).EnsureIndex(idx); err != nil {
			return nil, erro.Wrap(err)
		}
	}

	drv.sess = sess
	return drv.sess, nil
}

// drv.sess が oldSess なら drv.sess を newSess に変えて true を返す。
// そうでなければ false を返す。
func (drv *mongoDriver) replaceSession(newSess, oldSess *mgo.Session) bool {
	drv.sessLock.Lock()
	defer drv.sessLock.Unlock()
	if drv.sess == oldSess {
		drv.sess = newSess
		return true
	} else {
		return false
	}
}

func (drv *mongoDriver) closeIfError() error {
	drv.sessLock.Lock()
	sess := drv.sess
	drv.sessLock.Unlock()

	if sess == nil {
		return nil
	}

	if err := sess.Ping(); err == nil {
		// 問題無かった。
		return nil
	}

	if drv.replaceSession(nil, sess) {
		sess.Close()
	}
	return nil
}

func (drv *mongoDriver) close() error {
	drv.sessLock.Lock()
	defer drv.sessLock.Unlock()

	if drv.sess == nil {
		return nil
	}
	drv.sess.Close()
	drv.sess = nil
	return nil
}
