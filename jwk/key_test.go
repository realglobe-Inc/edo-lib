// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwk

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	cryptoutil "github.com/realglobe-Inc/edo-lib/crypto"
	"reflect"
	"testing"
)

func TestSample1(t *testing.T) {
	// JWK Appendix A.1 より。
	m := map[string]interface{}{
		"kty": "RSA",
		"n":   "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw",
		"e":   "AQAB",
	}

	if key, err := FromMap(m); err != nil {
		t.Error(err)
	} else if key.Type() != "RSA" {
		t.Error(key.Type())
		t.Error("RSA")
	} else if _, ok := key.Public().(*rsa.PublicKey); !ok {
		t.Error(key.Public())
	} else if buff := key.ToMap(); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestSample2(t *testing.T) {
	// JWK Appendix A.1 より。
	m := map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y":   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
	}

	if key, err := FromMap(m); err != nil {
		t.Error(err)
	} else if key.Type() != "EC" {
		t.Error(key.Type())
		t.Error("EX")
	} else if _, ok := key.Public().(*ecdsa.PublicKey); !ok {
		t.Error(key.Public())
	} else if buff := key.ToMap(); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestSample3(t *testing.T) {
	// JWE Appendix A.2 より。

	encrypted := []byte{
		80, 104, 72, 58, 11, 130, 236, 139, 132, 189, 255, 205, 61, 86, 151, 176,
		99, 40, 44, 233, 176, 189, 205, 70, 202, 169, 72, 40, 226, 181, 156, 223,
		120, 156, 115, 232, 150, 209, 145, 133, 104, 112, 237, 156, 116, 250, 65, 102,
		212, 210, 103, 240, 177, 61, 93, 40, 71, 231, 223, 226, 240, 157, 15, 31,
		150, 89, 200, 215, 198, 203, 108, 70, 117, 66, 212, 238, 193, 205, 23, 161,
		169, 218, 243, 203, 128, 214, 127, 253, 215, 139, 43, 17, 135, 103, 179, 220,
		28, 2, 212, 206, 131, 158, 128, 66, 62, 240, 78, 186, 141, 125, 132, 227,
		60, 137, 43, 31, 152, 199, 54, 72, 34, 212, 115, 11, 152, 101, 70, 42,
		219, 233, 142, 66, 151, 250, 126, 146, 141, 216, 190, 73, 50, 177, 146, 5,
		52, 247, 28, 197, 21, 59, 170, 247, 181, 89, 131, 241, 169, 182, 246, 99,
		15, 36, 102, 166, 182, 172, 197, 136, 230, 120, 60, 58, 219, 243, 149, 94,
		222, 150, 154, 194, 110, 227, 225, 112, 39, 89, 233, 112, 207, 211, 241, 124,
		174, 69, 221, 179, 107, 196, 225, 127, 167, 112, 226, 12, 242, 16, 24, 28,
		120, 182, 244, 213, 244, 153, 194, 162, 69, 160, 244, 248, 63, 165, 141, 4,
		207, 249, 193, 79, 131, 0, 169, 233, 127, 167, 101, 151, 125, 56, 112, 111,
		248, 29, 232, 90, 29, 147, 110, 169, 146, 114, 165, 204, 71, 136, 41, 252,
	}
	plain := []byte{
		4, 211, 31, 197, 84, 157, 252, 254, 11, 100, 157, 250, 63, 170, 106, 206,
		107, 124, 212, 45, 111, 107, 9, 219, 200, 177, 0, 240, 143, 156, 44, 207,
	}
	m := map[string]interface{}{
		"kty": "RSA",
		"n":   "sXchDaQebHnPiGvyDOAT4saGEUetSyo9MKLOoWFsueri23bOdgWp4Dy1WlUzewbgBHod5pcM9H95GQRV3JDXboIRROSBigeC5yjU1hGzHHyXss8UDprecbAYxknTcQkhslANGRUZmdTOQ5qTRsLAt6BTYuyvVRdhS8exSZEy_c4gs_7svlJJQ4H9_NxsiIoLwAEk7-Q3UXERGYw_75IDrGA84-lA_-Ct4eTlXHBIY2EaV7t7LjJaynVJCpkv4LKjTTAumiGUIuQhrNhZLuF_RJLqHpM2kgWFLU7-VTdL1VbC2tejvcI2BlMkEpk1BzBZI0KQB0GaDWFLN-aEAw3vRw",
		"e":   "AQAB",
		"d":   "VFCWOqXr8nvZNyaaJLXdnNPXZKRaWCjkU5Q2egQQpTBMwhprMzWzpR8Sxq1OPThh_J6MUD8Z35wky9b8eEO0pwNS8xlh1lOFRRBoNqDIKVOku0aZb-rynq8cxjDTLZQ6Fz7jSjR1Klop-YKaUHc9GsEofQqYruPhzSA-QgajZGPbE_0ZaVDJHfyd7UUBUKunFMScbflYAAOYJqVIVwaYR5zWEEceUjNnTNo_CVSj-VvXLO5VZfCUAVLgW4dpf1SrtZjSt34YLsRarSb127reG_DUwg9Ch-KyvjT1SkHgUWRVGcyly7uvVGRSDwsXypdrNinPA4jlhoNdizK2zF2CWQ",
		"p":   "9gY2w6I6S6L0juEKsbeDAwpd9WMfgqFoeA9vEyEUuk4kLwBKcoe1x4HG68ik918hdDSE9vDQSccA3xXHOAFOPJ8R9EeIAbTi1VwBYnbTp87X-xcPWlEPkrdoUKW60tgs1aNd_Nnc9LEVVPMS390zbFxt8TN_biaBgelNgbC95sM",
		"q":   "uKlCKvKv_ZJMVcdIs5vVSU_6cPtYI1ljWytExV_skstvRSNi9r66jdd9-yBhVfuG4shsp2j7rGnIio901RBeHo6TPKWVVykPu1iYhQXw1jIABfw-MVsN-3bQ76WLdt2SDxsHs7q7zPyUyHXmps7ycZ5c72wGkUwNOjYelmkiNS0",
		"dp":  "w0kZbV63cVRvVX6yk3C8cMxo2qCM4Y8nsq1lmMSYhG4EcL6FWbX5h9yuvngs4iLEFk6eALoUS4vIWEwcL4txw9LsWH_zKI-hwoReoP77cOdSL4AVcraHawlkpyd2TWjE5evgbhWtOxnZee3cXJBkAi64Ik6jZxbvk-RR3pEhnCs",
		"dq":  "o_8V14SezckO6CNLKs_btPdFiO9_kC1DsuUTd2LAfIIVeMZ7jn1Gus_Ff7B7IVx3p5KuBGOVF8L-qifLb6nQnLysgHDh132NDioZkhH7mI7hPG-PYE_odApKdnqECHWw0J-F0JWnUd6D2B_1TvF9mXA2Qx-iGYn8OVV1Bsmp6qU",
		"qi":  "eNho5yRBEBxhGBtQRww9QirZsB66TrfFReG_CcteI1aCneT0ELGhYlRlCtUkTRclIfuEPmNsNDPbLoLqqCVznFbvdB7x-Tl-m0l_eFTj2KiqwGqE9PZB9nNTwMVvH3VRRSLWACvPnSiwP8N5Usy-WRXS-V7TbpxIhvepTfE0NNo",
	}

	if key, err := FromMap(m); err != nil {
		t.Error(err)
	} else if key.Type() != "RSA" {
		t.Error(key.Type())
		t.Error("RSA")
	} else if pri, ok := key.Private().(*rsa.PrivateKey); !ok {
		t.Error(key.Private())
	} else if buff := key.ToMap(); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	} else if decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, pri, encrypted); err != nil {
		t.Error(err)
	} else if !bytes.Equal(decrypted, plain) {
		t.Error(decrypted)
		t.Error(plain)
	}
}

func TestSample4(t *testing.T) {
	// JWK Appendix A.2 より。

	m := map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y":   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"d":   "870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
	}

	if key, err := FromMap(m); err != nil {
		t.Error(err)
	} else if key.Type() != "EC" {
		t.Error(key.Type())
		t.Error("EC")
	} else if _, ok := key.Private().(*ecdsa.PrivateKey); !ok {
		t.Error(key.Private())
	} else if buff := key.ToMap(); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

var test_128Key, test_192Key, test_256Key, test_384Key, test_512Key []byte
var test_rsaKey *rsa.PrivateKey
var test_ec256Key, test_ec384Key, test_ec521Key *ecdsa.PrivateKey

func init() {
	for ; len(test_512Key) < 64; test_512Key = append(test_512Key, byte(len(test_512Key))) {
	}
	test_128Key = test_512Key[:16]
	test_192Key = test_512Key[:24]
	test_256Key = test_512Key[:32]
	test_384Key = test_512Key[:48]

	if k, err := cryptoutil.ParsePem([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA2Tv/db3EjCJYgY4USmfVd2yjMZroIV4b9VR2SgW71Rq85lik
qd8hvcXI6BZV6rdKpqDEzQ1u+n3R5L2+9SCrece4HbngprXsvuYluLEtOkdrJy14
ne55vawEstF9TRCcw7GT0aY+x91Uc+3Lc+D9FirkZua4Bz+mqQfFoe37XXZr/FFx
k0yws9tyzbnj9Rl4lJs9fecVTP8qNjg8UBOTaVfb716Qud3fKepTkyhyeet3x5tA
n/MTGJAbGhoFnzFvcyUFbY9eGv67wd50TxXD9RAo98IfBNLxXukxngnFeCSRom4u
1Pvh4UNqedOBOF6hU1CmRdAzksvd6zqhRCpjQQIDAQABAoIBAQC11WPS4WKIziL5
Zr0TPvDOww+i8QBHFegfJXDSKxR7n6LoyOAkFNLAb7LomfGWw4/oBABXh1wSroin
iDA0LQF7sTIrJ7CkuvkNHcYLX7r04l0N8SDaSYh7vGY+a94PSM1/fL+3qAk68MfF
NhGr0HLoQETo4Uy/PIc7S3chQPu70V1ZI4hm7Q4T9+mgQ7O+Y6c07KOhw5bElGNM
8VUUskTiZERfGSxvcdnSDY2W8Wq6IY1G1YPjN0wnj1TKMSA7d4lD59auzI8k+LTx
L3XvXmOOB3UygTPPbGsW4hs7sg+GJWaOB1U8RHrUsyDtOIocKSWFTsyVHwfulKce
RipUxwsBAoGBAP3v69ST/RXqmMhXk24WaAw1Jv72s4hQLmcCw38FrbT527FsxP6b
gxHjmPD27LyXe35m19raZ0qDn7mFum934d0saRzXLwxyu6mwttgxiP8JMWNNWh1f
y30vRpYwqlz+2wVnaFwKnPOu/CZkENjCcM51ZKeV9Tel7x3qCaAoPAeRAoGBANr/
wDNH6IlN/90/h6c/YoHLfjp097ZL1LmeZ2zeno1lQUTdqILQSzf0lX4MTICzU3MR
tL7XVaTGaH7EaGxZ8gPbRh+dKRRrK9OFH94Dr4W39UAelgcaC4aFIesTCdWNQG9l
1w1+mWrgDm10FyfzeIYHJCGAL+Tg2iqEa/ltkqixAoGAZwZY+sUT0Dl+xQFq6iYj
Dpjd+mFi03IccWSYpkdKg3s/m8tSXS4Azlg1q8WypI0c6FqXRs6HS579RYqw6hqM
Q2yKNM5E41sFMkJk3G+0cixrois23WYJK//rNnIGHHa1q4qZt4YCyYb7/CNrBlZU
6B6OuMNJWstyqQNT5muMd1ECgYBzvW6Kq5pN1pc/CvBah7k795wCsQappXILl5f5
hb4t5DGWf78rQ4I9VFodf8p+ykd0LQtlQNDWgLWBKbQ2b2LkfuKUmq63R9ylsVmi
MDh3Zz2KYZ/QqQcmVP4UCr/LyRcgyKXbT+ks/rUhS5VhW996lhOWUPT+9YbXqZyW
+j3kQQKBgQDy0qmVxszmSglkSdmuO6UNA2Z7UJAKUesSKn6c+u9JBJ3zIALkxUUu
XUZirwxI0tP14ntYvaf0w6PVsHLbrj2piITQbZcRQuBg/JW7mByI4KidcONJ+Fxl
bFlw4ybcrNNPrjnVqfqt4cNh59yN2WAJCZTj3xOHLtdmcmELBQVeAQ==
-----END RSA PRIVATE KEY-----`)); err != nil {
		panic(err)
	} else {
		test_rsaKey = k.(*rsa.PrivateKey)
	}

	if k, err := cryptoutil.ParsePem([]byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPFUo1nmauOJltxl0nfaVx3BEZ6wdg+hRI+S8OfUIDQaoAoGCCqGSM49
AwEHoUQDQgAE3tfF/QYgrjnyDzRPycEyx0yZUvX2xZS8JFQb74c91Oi5OtThEZDq
iyltctMoRBmc1JBq9Doh5ZybUQio1aV46A==
-----END EC PRIVATE KEY-----`)); err != nil {
		panic(err)
	} else {
		test_ec256Key = k.(*ecdsa.PrivateKey)
	}

	if k, err := cryptoutil.ParsePem([]byte(`-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDB15rTMYJQzyMcZWBEpciIj7R4suAEyh+bvX6VYbLbPpdLchP09WtmJ
DoGNuRrYqrugBwYFK4EEACKhZANiAAQcvQFWzpz0PF7PuHNpj1GGZWo5w8hxUMyf
kyeiq/+zNK51qNYRCQSoshBOv5uqvABmVF0Sgg3uFsJEoVWQMPezRvIWPhE37zx5
ux83IWcL9esbhYcCuocvoOWi3CiYBHY=
-----END EC PRIVATE KEY-----`)); err != nil {
		panic(err)
	} else {
		test_ec384Key = k.(*ecdsa.PrivateKey)
	}

	if k, err := cryptoutil.ParsePem([]byte(`-----BEGIN EC PRIVATE KEY-----
MIHbAgEBBEEe5zLkeThtmXAllBN08a/rhTtIEuXlGwA10gaRaKpDCJ+sg5RKAGiE
qyd6830JEvDRqOFOoN/6CmHEL3FWd1dp5qAHBgUrgQQAI6GBiQOBhgAEADeJGCqZ
XlevrSbI8QTfYtEFWAVqgXJKfPLxZO1RX8MY834ApWIkj2Dkx61Y5LUXuKmRvUqb
4PaPe5YUgS3ysAWdAFK8IykafrwtCjdmF/hQy/C2KlYNk9RYh162KSPlyJIZV7Xb
pTJ/8vECYoRmB0BZ9c9RmTILbjqBzd8+Hb3TvR0C
-----END EC PRIVATE KEY-----`)); err != nil {
		panic(err)
	} else {
		test_ec521Key = k.(*ecdsa.PrivateKey)
	}
}

func TestKey(t *testing.T) {
	if key := New(test_128Key, map[string]interface{}{
		"use": "sig",
	}); key.Use() != "sig" {
		t.Error(key.Use())
		t.Error("sig")
	}

	if key := New(test_128Key, map[string]interface{}{
		"key_ops": []interface{}{"sign", "verify"},
	}); !reflect.DeepEqual(key.Operations(), map[string]bool{"sign": true, "verify": true}) {
		t.Error(key.Operations())
		t.Error(map[string]bool{"sign": true, "verify": true})
	}

	if key := New(test_128Key, map[string]interface{}{
		"alg": "A128KW",
	}); key.Algorithm() != "A128KW" {
		t.Error(key.Algorithm())
		t.Error("A128KW")
	}

	if key := New(test_128Key, map[string]interface{}{
		"kid": "test-key",
	}); key.Id() != "test-key" {
		t.Error(key.Id())
		t.Error("test-key")
	}
}

func TestToFromMap(t *testing.T) {
	for _, rawKey := range []interface{}{
		test_128Key,
		test_192Key,
		test_256Key,
		test_384Key,
		test_512Key,
		test_rsaKey,
		test_ec256Key,
		test_ec384Key,
		test_ec521Key,
		&test_rsaKey.PublicKey,
		&test_ec256Key.PublicKey,
		&test_ec384Key.PublicKey,
		&test_ec521Key.PublicKey,
	} {
		key := New(rawKey, nil)
		if key2, err := FromMap(key.ToMap()); err != nil {
			t.Error(err)
			t.Error(key)
		} else if !reflect.DeepEqual(key2, key) {
			t.Error(key2)
			t.Error(key)
		}
	}
}
