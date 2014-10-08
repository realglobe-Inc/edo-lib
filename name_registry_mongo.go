package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

// mongodb から map[string]string を読み取る。
func containerMongoTake(query *mgo.Query) (interface{}, *Stamp, error) {
	var res struct {
		Value map[string]string
		Stamp *Stamp
	}
	if err := query.One(&res); err != nil {
		return nil, nil, erro.Wrap(err)
	}
	return res.Value, res.Stamp, nil
}

// スレッドセーフ。
func NewMongoNameRegistry(url, dbName, collName string, expiDur time.Duration) (NameRegistry, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.MongoMarshal = func(value interface{}) (interface{}, error) {
		cont := map[string]string{}
		for k, v := range value.(*nameTree).toContainer() {
			cont[strings.Replace(k, ".", "/", -1)] = v
		}
		return cont, nil
	}
	base.MongoUnmarshal = func(data interface{}) (interface{}, error) {
		tree := newNameTree()
		for key, value := range data.(map[string]string) {
			tree.add(strings.Replace(key, "/", ".", -1), value)
		}
		return tree, nil
	}
	base.MongoTake = containerMongoTake
	return newNameRegistry(base), nil
}
