package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"strconv"
	"time"
)

// 非キャッシュ用。
func NewMongoJsRegistry(url, dbName, collName string) (JsRegistry, error) {
	return newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"path"},
			Unique:   true,
			DropDups: true,
		},
	})
}

type mongoObject struct {
	Path string `bson:"path"`

	Service bool     `bson:"service,omitempty"`
	Library bool     `bson:"library,omitempty"`
	Include []string `bson:"include,omitempty"`
	Code    string   `bson:"code"`

	Date   time.Time `bson:"date"`
	Digest int       `bson:"digest"`
}

func (reg *mongoDriver) Object(dir, objName string) (*Object, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"path": path.Join(dir, objName)})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	var res mongoObject
	if err := query.One(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &Object{res.Service, res.Library, res.Include, res.Code}, nil
}

func (obj *Object) digest() int {
	prime := 31
	dig := 0
	dig = prime*dig + util.DigestBool(obj.Service)
	dig = prime*dig + util.DigestBool(obj.Library)
	for _, inc := range obj.Include {
		dig = prime*dig + util.DigestString(inc)
	}
	dig = prime*dig + util.DigestString(obj.Code)
	return dig

}

func (reg *mongoDriver) AddObject(dir, objName string, obj *Object) error {
	mongoObj := &mongoObject{path.Join(dir, objName), obj.Service, obj.Library, obj.Include, obj.Code, time.Now(), obj.digest()}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"path": mongoObj.Path}, mongoObj); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoDriver) RemoveObject(dir, objName string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{"path": path.Join(dir, objName)}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// キャッシュ用。
func NewMongoJsBackendRegistry(url, dbName, collName string, expiDur time.Duration) (JsBackendRegistry, error) {
	reg, err := NewMongoJsRegistry(url, dbName, collName)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return newDatedMongoDriver(reg.(*mongoDriver), expiDur), nil
}

func (reg *datedMongoDriver) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"path": path.Join(dir, objName)})
	if n, err := query.Count(); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil, nil
	}
	var res mongoObject
	if err := query.One(&res); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	stmp := &Stamp{Date: res.Date, Digest: strconv.Itoa(res.Digest)}

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	return &Object{res.Service, res.Library, res.Include, res.Code}, newCaStmp, nil
}
