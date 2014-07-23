package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// mondodb を使うドライバー。

// ジョブ。
type mongoRegistry struct {
	dbName   string
	collName string
	*mgo.Session
}

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
