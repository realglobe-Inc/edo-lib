package driver

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"
	"time"
)

// 事前に、サービス UUID a_b-c、公開鍵 testPublicKey で登録しとく。

var testPublicKey *rsa.PublicKey
var testPublicKeyPem string

func init() {
	testKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	testPublicKey = &testKey.PublicKey
	testPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(testPublicKey)
	testPublicKeyPem = string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: testPublicKeyBytes,
		},
	))
}

func TestParseKey(t *testing.T) {
	if key, err := parseKey(testPublicKeyPem); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(key, testPublicKey) {
		t.Error(key, testPublicKey)
	}
}

// 非キャッシュ用。
func testServiceKeyRegistry(t *testing.T, reg ServiceKeyRegistry) {
	key, err := reg.ServiceKey(testServUuid)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(key, testPublicKey) {
		t.Error(key, testPublicKey)
	}
}

// キャッシュ用。
func testDatedServiceKeyRegistry(t *testing.T, reg DatedServiceKeyRegistry) {

	key1, stmp1, err := reg.StampedServiceKey(testServUuid, nil)
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(key1, testPublicKey) || stmp1 == nil {
		t.Error(key1, stmp1)
	}

	// キャッシュと同じだから返らない。
	key2, stmp2, err := reg.StampedServiceKey(testServUuid, stmp1)
	if err != nil {
		t.Fatal(err)
	} else if key2 != nil || stmp2 == nil {
		t.Error(key2, stmp2)
	}

	// キャッシュが古いから返る。
	key3, stmp3, err := reg.StampedServiceKey(testServUuid, &Stamp{Date: stmp1.Date.Add(-time.Second), Digest: stmp1.Digest})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(key3, testPublicKey) || stmp3 == nil {
		t.Error(key3, stmp3)
	}

	// ダイジェストが違うから返る。
	key4, stmp4, err := reg.StampedServiceKey(testServUuid, &Stamp{Date: stmp1.Date, Digest: stmp1.Digest + "a"})
	if err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(key3, testPublicKey) || stmp4 == nil {
		t.Error(key4, stmp4)
	}
}
