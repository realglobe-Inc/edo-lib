package jwt

import (
	"reflect"
	"testing"
)

func TestJwkRsa(t *testing.T) {
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

func TestJwkEcdsa(t *testing.T) {
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
