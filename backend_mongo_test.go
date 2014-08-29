package driver

import (
	"testing"
	"time"
)

func TestMongoJsBackendRegistry(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoJsBackendRegistry(mongoAddr, "test_driver_mongo", "js")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo").DropDatabase()

	testJsBackendRegistry(t, reg)
}

func TestMongoIdProviderBackend(t *testing.T) {
	if mongoAddr == "" {
		t.SkipNow()
	}

	reg, err := NewMongoIdProviderBackend(mongoAddr, "test_driver_mongo_id_provider_backend", "idp")
	if err != nil {
		t.Fatal(err)
	}
	defer reg.(*mongoRegistry).DB("test_driver_mongo_id_provider_backend").DropDatabase()

	if err := reg.(*mongoRegistry).DB("test_driver_mongo_id_provider_backend").C("idp").Insert(
		&IdProvider{
			IdpUuid: "a_b-c",
			Name:    "ABC",
			Uri:     "https://localhost:1234",
		},
	); err != nil {
		t.Fatal(err)
	}
	if err := reg.(*mongoRegistry).DB("test_driver_mongo_id_provider_backend").C("idp").Insert(
		&Stamp{
			Date:   time.Now(),
			Digest: "0",
		},
	); err != nil {
		t.Fatal(err)
	}

	testIdProviderBackend(t, reg)
}
