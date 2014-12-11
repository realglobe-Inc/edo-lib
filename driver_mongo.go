package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"reflect"
	"sync"
)

// スレッドセーフ。
type mongoDriver struct {
	url     string
	db      string
	coll    string
	indices []mgo.Index

	sLock sync.Mutex
	s     *mgo.Session
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
	reg.sLock.Lock()
	defer reg.sLock.Unlock()

	if reg.s != nil {
		return reg.s, nil
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

	reg.s = sess
	return reg.s, nil
}

// reg.s が oldSess なら reg.s を newSess に変えて true を返す。
// そうでなければ false を返す。
func (reg *mongoDriver) replaceSession(newSess, oldSess *mgo.Session) bool {
	reg.sLock.Lock()
	defer reg.sLock.Unlock()
	if reg.s == oldSess {
		reg.s = newSess
		return true
	} else {
		return false
	}
}

func (reg *mongoDriver) closeIfError() error {
	reg.sLock.Lock()
	sess := reg.s
	reg.sLock.Unlock()

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
