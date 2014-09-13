package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 非キャッシュ用。
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
