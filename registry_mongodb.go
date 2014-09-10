package driver

import (
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"path"
	"time"
)

// ログイン。
func NewMongoLoginRegistry(url, dbName, collName string) (LoginRegistry, error) {
	return newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"access_token"},
			Unique:   true,
			DropDups: true,
		},
	})
}

type mongoUser struct {
	AccToken string `bson:"access_token"`

	UsrUuid string `bson:"user_uuid"`
}

func (reg *mongoDriver) User(accToken string) (usrUuid string, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"access_token": accToken})
	if n, err := query.Count(); err != nil {
		return "", erro.Wrap(err)
	} else if n == 0 {
		return "", nil
	}
	var res mongoUser
	if err := query.One(&res); err != nil {
		return "", erro.Wrap(err)
	}
	return res.UsrUuid, nil
}

// JavaScript.
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

// ユーザー情報。
func NewMongoUserRegistry(url, dbName, collName string) (UserRegistry, error) {
	return newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key: []string{"user_uuid"},
		},
		mgo.Index{
			Key:      []string{"user_uuid", "key"},
			Unique:   true,
			DropDups: true,
		},
	})
}

type mongoAttribute struct {
	UsrUuid  string      `bson:"user_uuid"`
	AttrName string      `bson:"key"`
	Attr     interface{} `bson:"value"`
}

func (reg *mongoDriver) Attributes(usrUuid string) (attrs map[string]interface{}, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"user_uuid": usrUuid})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	mongoAttrs := []mongoAttribute{}
	if err := query.Iter().All(&mongoAttrs); err != nil {
		return nil, erro.Wrap(err)
	}
	attrs = map[string]interface{}{}
	for _, mongoAttr := range mongoAttrs {
		attrs[mongoAttr.AttrName] = mongoAttr.Attr
	}
	return attrs, nil
}

func (reg *mongoDriver) Attribute(usrUuid, attrName string) (attr interface{}, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"user_uuid": usrUuid, "key": attrName})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	var mongoAttr mongoAttribute
	if err := query.One(&mongoAttr); err != nil {
		return nil, erro.Wrap(err)
	}
	return mongoAttr.Attr, nil
}

func (reg *mongoDriver) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	mongoAttr := &mongoAttribute{usrUuid, attrName, attr}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"user_uuid": usrUuid, "key": attrName}, mongoAttr); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoDriver) RemoveAttribute(usrUuid, attrName string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{"user_uuid": usrUuid, "key": attrName}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// ジョブ。
func NewMongoJobRegistry(url, dbName, collName string) (JobRegistry, error) {
	return newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"job_id"},
			Unique:   true,
			DropDups: true,
		},
		mgo.Index{
			Key: []string{"deadline"},
		},
	})
}

type mongoJobResult struct {
	JobId    string            `bson:"job_id"`
	Deadline time.Time         `bson:"deadline"`
	Status   int               `bson:"status"`
	Headers  map[string]string `bson:"headers,omitempty"`
	Body     string            `bson:"body,omitempty"`
}

func (reg *mongoDriver) Result(jobId string) (*JobResult, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"job_id": jobId})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	var res mongoJobResult
	if err := query.One(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &JobResult{res.Status, res.Headers, res.Body}, nil
}

func (reg *mongoDriver) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	mongoRes := &mongoJobResult{jobId, deadline, res.Status, res.Headers, res.Body}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"job_id": jobId}, mongoRes); err != nil {
		return erro.Wrap(err)
	}

	// 古いのを削除。
	reg.DB(reg.dbName).C(reg.collName).RemoveAll(bson.M{"deadline": bson.M{"$lt": time.Now()}})

	return nil
}

// 別名。
func NewMongoNameRegistry(url, dbName, collName string) (NameRegistry, error) {
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	nameIdx := mgo.Index{
		Key:      []string{"name"},
		Unique:   true,
		DropDups: true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(nameIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	return &mongoDriver{dbName, collName, sess}, nil
}

type mongoAddress struct {
	Name string `bson:"name"`
	Addr string `bson:"address"`
}

func (reg *mongoDriver) Address(name string) (addr string, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"name": name})
	if n, err := query.Count(); err != nil {
		return "", erro.Wrap(err)
	} else if n == 0 {
		return "", nil
	}
	var mongoAddr mongoAddress
	if err := query.One(&mongoAddr); err != nil {
		return "", erro.Wrap(err)
	}
	return mongoAddr.Addr, nil
}

func (reg *mongoDriver) Addresses(name string) (addrs []string, err error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"name": bson.M{"$regex": name + "$"}})
	var mongoAddrs []mongoAddress
	if err := query.All(&mongoAddrs); err != nil {
		return nil, erro.Wrap(err)
	}
	for _, mongoAddr := range mongoAddrs {
		addrs = append(addrs, mongoAddr.Addr)
	}
	return addrs, nil
}

// イベント。
func NewMongoEventRegistry(url, dbName, collName string) (EventRegistry, error) {
	return newMongoDriver(url, dbName, collName, []mgo.Index{
		mgo.Index{
			Key:      []string{"user_uuid", "event"},
			Unique:   true,
			DropDups: true,
		},
	})
}

type mongoHandler struct {
	UsrUuid string  `bson:"user_uuid"`
	Event   string  `bson:"event"`
	Hndl    Handler `bson:"handler"`
}

func (reg *mongoDriver) Handler(usrUuid, event string) (Handler, error) {
	query := reg.DB(reg.dbName).C(reg.collName).Find(bson.M{"user_uuid": usrUuid, "event": event})
	if n, err := query.Count(); err != nil {
		return nil, erro.Wrap(err)
	} else if n == 0 {
		return nil, nil
	}
	var res mongoHandler
	if err := query.One(&res); err != nil {
		return nil, erro.Wrap(err)
	}
	return res.Hndl, nil
}

func (reg *mongoDriver) AddHandler(usrUuid, event string, hndl Handler) error {
	mongoHndl := &mongoHandler{usrUuid, event, hndl}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"user_uuid": usrUuid, "event": event}, mongoHndl); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoDriver) RemoveHandler(usrUuid, event string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{"user_uuid": usrUuid, "event": event}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
