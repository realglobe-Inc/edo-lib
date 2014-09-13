package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 非キャッシュ用。
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
