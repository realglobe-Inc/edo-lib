package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 非キャッシュ用。
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
