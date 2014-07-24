package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// mondodb を使うドライバー。

type mongoRegistry struct {
	dbName   string
	collName string
	*mgo.Session
}

// ユーザー情報。
func NewMongoUserRegistry(url, dbName, collName string) (UserRegistry, error) {
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	usrUuidIdx := mgo.Index{
		Key:      []string{"user_uuid"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(usrUuidIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	keyIdx := mgo.Index{
		Key:      []string{"user_uuid", "key"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(keyIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	return &mongoRegistry{dbName, collName, sess}, nil
}

type mongoAttribute struct {
	UsrUuid  string      `bson:"user_uuid"`
	AttrName string      `bson:"key"`
	Attr     interface{} `bson:"value"`
}

func (reg *mongoRegistry) Attributes(usrUuid string) (attrs map[string]interface{}, err error) {
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

func (reg *mongoRegistry) Attribute(usrUuid, attrName string) (attr interface{}, err error) {
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

func (reg *mongoRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	mongoAttr := &mongoAttribute{usrUuid, attrName, attr}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"user_uuid": usrUuid, "key": attrName}, mongoAttr); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoRegistry) RemoveAttribute(usrUuid, attrName string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{"user_uuid": usrUuid, "key": attrName}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// ジョブ。
func NewMongoJobRegistry(url, dbName, collName string) (JobRegistry, error) {
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	jobIdIdx := mgo.Index{
		Key:      []string{"job_id"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(jobIdIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	deadlineIdx := mgo.Index{
		Key: []string{"deadline"},
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(deadlineIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	return &mongoRegistry{dbName, collName, sess}, nil
}

type mongoJobResult struct {
	JobId    string            `bson:"job_id"`
	Deadline time.Time         `bson:"deadline"`
	Status   int               `bson:"status"`
	Headers  map[string]string `bson:"headers,omitempty"`
	Body     string            `bson:"body,omitempty"`
}

func (reg *mongoRegistry) Result(jobId string) (*JobResult, error) {
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

func (reg *mongoRegistry) AddResult(jobId string, res *JobResult, deadline time.Time) error {
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

	jobIdIdx := mgo.Index{
		Key:      []string{"name"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(jobIdIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	return &mongoRegistry{dbName, collName, sess}, nil
}

type mongoAddress struct {
	Name string `bson:"name"`
	Addr string `bson:"address"`
}

func (reg *mongoRegistry) Address(name string) (addr string, err error) {
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

func (reg *mongoRegistry) Addresses(name string) (addrs []string, err error) {
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
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	eventIdx := mgo.Index{
		Key:      []string{"user_uuid", "event"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}
	if err := sess.DB(dbName).C(collName).EnsureIndex(eventIdx); err != nil {
		return nil, erro.Wrap(err)
	}

	return &mongoRegistry{dbName, collName, sess}, nil
}

type mongoHandler struct {
	UsrUuid string  `bson:"user_uuid"`
	Event   string  `bson:"event"`
	Hndl    Handler `bson:"handler"`
}

func (reg *mongoRegistry) Handler(usrUuid, event string) (Handler, error) {
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

func (reg *mongoRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	mongoHndl := &mongoHandler{usrUuid, event, hndl}
	if _, err := reg.DB(reg.dbName).C(reg.collName).Upsert(bson.M{"user_uuid": usrUuid, "event": event}, mongoHndl); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func (reg *mongoRegistry) RemoveHandler(usrUuid, event string) error {
	if err := reg.DB(reg.dbName).C(reg.collName).Remove(bson.M{"user_uuid": usrUuid, "event": event}); err != nil {
		return erro.Wrap(err)
	}
	return nil
}
