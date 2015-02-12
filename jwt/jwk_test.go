package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
)

var testRsaKey *rsa.PrivateKey
var testEcdsaKey *ecdsa.PrivateKey
var testEcdsa256Key *ecdsa.PrivateKey
var testEcdsa384Key *ecdsa.PrivateKey
var testEcdsa521Key *ecdsa.PrivateKey

func init() {
	var err error
	testRsaKey, err = rsa.GenerateKey(rand.Reader, 1152)
	if err != nil {
		panic(err)
	}
	testEcdsa256Key, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	testEcdsa384Key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		panic(err)
	}
	testEcdsa521Key, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	testEcdsaKey = testEcdsa256Key
}

func TestRsaPublickKey(t *testing.T) {
	m := map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	}

	_, key, err := PublicKeyFromJwkMap(m)
	if err != nil {
		t.Fatal(err)
	}

	if buff := PublicKeyToJwkMap("", key); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestEcdsaPublicKey(t *testing.T) {
	m := map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y":   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
	}

	_, key, err := PublicKeyFromJwkMap(m)
	if err != nil {
		t.Fatal(err)
	}

	if buff := PublicKeyToJwkMap("", key); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestPublicKey(t *testing.T) {
	for _, key := range []crypto.PublicKey{&testRsaKey.PublicKey, &testEcdsaKey.PublicKey} {
		_, key2, err := PublicKeyFromJwkMap(PublicKeyToJwkMap("", key))
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(key2, key) {
			t.Error(key2)
			t.Error(key)
		}
	}
}
