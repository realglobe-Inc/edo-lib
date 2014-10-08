package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"time"
)

// スレッドセーフ。
func NewMongoJsRegistry(url, dbName, collName string, expiDur time.Duration) (JsRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.MongoTake = func(query *mgo.Query) (interface{}, *Stamp, error) {
		var res struct {
			Value *Object
			Stamp *Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return res.Value, res.Stamp, nil
	}
	return newJsRegistry(base), nil
}
