package driver

import (
	"gopkg.in/mgo.v2"
	"testing"
)

// テストするなら、ローカルにデフォルトポートで mongodb をたてる必要あり。
var mongoAddr = "localhost"

func init() {
	if mongoAddr != "" {
		// 実際にサーバーが立っているかどうか調べる。
		// 立ってなかったらテストはスキップ。
		conn, err := mgo.Dial(mongoAddr)
		if err != nil {
			mongoAddr = ""
		} else {
			conn.Close()
		}
	}
}

func TestMongoLoginRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoLoginRegistry(mongoAddr, "test_driver_mongo", "login")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("login").Insert(
		&mongoUser{"abc-012", "a_b-c"},
	); err != nil {
		t.Fatal(err)
	}

	testLoginRegistry(t, reg)
}

func TestMongoJsRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJsRegistry(t, reg)
}

func TestMongoUserRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoUserRegistry(mongoAddr, "test_driver_mongo", "user")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testUserRegistry(t, reg)
}

func TestMongoJobRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJobRegistry(t, reg)
}

func TestMongoJobRegistryRemoveOld(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJobRegistry(mongoAddr, "test_driver_mongo", "job")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJobRegistryRemoveOld(t, reg)
}

func TestMongoNameRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoNameRegistry(mongoAddr, "test_driver_mongo", "name")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("name").Insert(
		&mongoAddress{"c.b.a", "c.localhost"},
		&mongoAddress{"d.b.a", "d.localhost"},
		&mongoAddress{"b.a", "localhost"},
	); err != nil {
		t.Fatal(err)
	}

	testNameRegistry(t, reg)
}

func TestMongoEventRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoEventRegistry(mongoAddr, "test_driver_mongo", "event")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testEventRegistry(t, reg)
}

func TestMongoServiceRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoServiceRegistry(mongoAddr, "test_driver_mongo", "service")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("service").Insert(
		&mongoService{map[string]string{
			"localhost:1234": "a_b-c",
		}},
	); err != nil {
		t.Fatal(err)
	}

	testServiceRegistry(t, reg)
}

func TestMongoLargeServiceRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoLargeServiceRegistry(mongoAddr, "test_driver_mongo", "service")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoLargeServiceRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoLargeServiceRegistry).DB("test_driver_mongo").C("service").Insert(
		&mongoLargeService{
			EndPt:    "localhost:1234",
			ServUuid: "a_b-c",
		},
	); err != nil {
		t.Fatal(err)
	}

	testServiceRegistry(t, reg)
}

func TestMongoIdProviderLister(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderLister(mongoAddr, "test_driver_mongo", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo").C("idp").Insert(
		&IdProvider{
			Uuid: "a_b-c",
			Name: "ABC",
			Uri:  "https://localhost:1234",
		},
	); err != nil {
		t.Fatal(err)
	}

	testIdProviderLister(t, reg)
}
