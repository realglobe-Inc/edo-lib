package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 非キャッシュ用。
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
