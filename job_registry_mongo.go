package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"time"
)

// スレッドセーフ。
func NewMongoJobRegistry(url, dbName, collName string, expiDur time.Duration) (JobRegistry, error) {
	base, err := newMongoTimeLimitedKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.base.MongoTake = func(query *mgo.Query) (interface{}, *Stamp, error) {
		var res struct {
			Value *JobResult
			Stamp *Stamp
		}
		if err := query.One(&res); err != nil {
			return nil, nil, erro.Wrap(err)
		}
		return res.Value, res.Stamp, nil
	}
	return newJobRegistry(base), nil
}
