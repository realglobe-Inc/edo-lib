package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
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

func (reg *mongoDriver) collection() (*mgo.Collection, error) {
	sess, err := reg.session()
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return sess.DB(reg.db).C(reg.coll), nil
}

func (reg *mongoDriver) session() (*mgo.Session, error) {
	reg.sessLock.Lock()
	defer reg.sessLock.Unlock()

	if reg.sess != nil {
		return reg.sess, nil
	}

	sess, err := mgo.Dial(reg.url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	curIndices, err := sess.DB(reg.db).C(reg.coll).Indexes()
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, idx := range reg.indices {
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

		if err := sess.DB(reg.db).C(reg.coll).EnsureIndex(idx); err != nil {
			return nil, erro.Wrap(err)
		}
	}

	reg.sess = sess
	return reg.sess, nil
}

// reg.sess が oldSess なら reg.sess を newSess に変えて true を返す。
// そうでなければ false を返す。
func (reg *mongoDriver) replaceSession(newSess, oldSess *mgo.Session) bool {
	reg.sessLock.Lock()
	defer reg.sessLock.Unlock()
	if reg.sess == oldSess {
		reg.sess = newSess
		return true
	} else {
		return false
	}
}

func (reg *mongoDriver) closeIfError() error {
	reg.sessLock.Lock()
	sess := reg.sess
	reg.sessLock.Unlock()

	if sess == nil {
		return nil
	}

	if err := sess.Ping(); err == nil {
		// 問題無かった。
		return nil
	}

	if reg.replaceSession(nil, sess) {
		sess.Close()
	}
	return nil
}
