package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 非キャッシュ用。
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
