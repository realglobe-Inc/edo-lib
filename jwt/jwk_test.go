package jwt

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"
)

var test128Key, test192Key, test256Key, test384Key, test512Key []byte
var testRsaKey *rsa.PrivateKey
var testEcdsaKey *ecdsa.PrivateKey
var testEcdsa256Key, testEcdsa384Key, testEcdsa521Key *ecdsa.PrivateKey

var testSigKeys, testVeriKeys map[string]interface{}
var testEncKeys, testDecKeys map[string]interface{}

func init() {
	for ; len(test512Key) < 64; test512Key = append(test512Key, byte(len(test512Key))) {
	}
	test128Key = test512Key[:16]
	test192Key = test512Key[:24]
	test256Key = test512Key[:32]
	test384Key = test512Key[:48]

	key, err := KeyFromJwkMap(map[string]interface{}{
		"d":   "gOV1-Oo5UenUbuT6xXWmsHOlCOriHaH-iis22HdliQAjMxaO0_Yog8pSG4bRit7xIn-_olkmRZm2X21gd2AUC_mkE7Nytw5t_pioMzupEVVGApIFuc2_ryf5VPSznx3zk5FY6XCgUf6BnJ188WRUv3CnnNuAmEJtP6MhWmoKlPMpgQ",
		"dp":  "sa3EMdvoT87Z-ecMyWpw-_EA-AICiynWHcaW8iYbc9r2inlfmJ-61mzRzOXITFA8x2nKOqOkT4eFYKIauHzaQ_U",
		"dq":  "EhoQ2ioI0VbueHV34SHmKSZIbkXTjuJD4hTzAEz-i6Wuma4lYpNxz3pI-mYXrVWdmjy07ErOou-vcuZ3gFMg_iE",
		"e":   "AQAB",
		"kty": "RSA",
		"n":   "5cadP6Vvv6ABglXpSeXYxPB321gtSwmjccsHr2-YKmBm22KWF2A1b68LJ3mA8eG5NPSRL6macCMttxsoAKwaCxOxn-6dNOKXNLQ1S0WsE4yY2QLoi9Cj_sY8yfdk_wb0ZM5kyE99GjFFLDvnh-RjHIf2cbXPyPfbeLigeeon7jsxOw",
		"p":   "9556qFgzilKEEhQ41fVzvLm5vKpiCc0IABG1CDQ_VTr4KGoOcqSHx6__yqYFQlzgizkG-zVxBQSs-6GZ3eA-t4s",
		"q":   "7Y2H2tRgIm9UjN0OlszOBcOXqPicE5KlseuNCIJZo1SyW30h-N2ssjCeiSDPrqm5QGZ637EAmhvNsPNOxzwLIxE",
		"qi":  "DbisQAteFbdCaNy6TyNy5UgZjdPba1bhKI3iIXalno_5HRrK4tUzu9VHdYVj5-iscIw5za9cPMLFr3zQvWa-gzA",
	})
	if err != nil {
		panic(err)
	}
	testRsaKey = key.(*rsa.PrivateKey)

	key, err = KeyFromJwkMap(map[string]interface{}{
		"crv": "P-256",
		"d":   "3BhkCluOkm8d8gvaPD5FDG2zeEw2JKf3D5LwN-mYmsw",
		"kty": "EC",
		"x":   "lpHYO1qpjU95B2sThPR2-1jv44axgaEDkQtcKNE-oZs",
		"y":   "soy5O11SFFFeYdhQVodXlYPIpeo0pCS69IxiVPPf0Tk",
	})
	if err != nil {
		panic(err)
	}
	testEcdsa256Key = key.(*ecdsa.PrivateKey)

	key, err = KeyFromJwkMap(map[string]interface{}{
		"crv": "P-384",
		"d":   "Gp-7eC0G7PjGzKoiAmTQ1iLsLU3AEy3h-bKFWSZOanXqSWI6wqJVPEUsatNYBJoG",
		"kty": "EC",
		"x":   "HlrMhzZww_AkmHV-2gDR5n7t75673UClnC7V2GewWva_sg-4GSUguFalVgwnK0tQ",
		"y":   "fxS48Fy50SZFZ-RAQRWUZXZgRSWwiKVkqPTd6gypfpQNkXSwE69BXYIAQcfaLcf2",
	})
	if err != nil {
		panic(err)
	}
	testEcdsa384Key = key.(*ecdsa.PrivateKey)

	key, err = KeyFromJwkMap(map[string]interface{}{
		"crv": "P-521",
		"d":   "AK9ejP0HuUt7ojjI9p20986DGqG-5jc9UWMtnMqxNvIvBScTJflS2lE-6CRsJZKk6ChWI6U4ahXDH0cCCFWTAvI_",
		"kty": "EC",
		"x":   "AXA9Y2pY1g_Cs4W6Nto7ebjDKOsaRHxo6EYRWjk1XHZaA7HnkeHg13x24OWHelqdZiuo7J1VbRKJ4ohZPjKX-AL7",
		"y":   "AeC5zdcUHLDdhQAvdsnnwD8rgNWjMdlWsqXZxv7Ar7ly5xmZGxEDtcJuhfhn8R9PXeScPH2soF3dFYPCuDkF4Gns",
	})
	if err != nil {
		panic(err)
	}
	testEcdsa521Key = key.(*ecdsa.PrivateKey)

	testEcdsaKey = testEcdsa256Key

	testSigKeys = map[string]interface{}{
		"none":  nil,
		"HS256": test256Key,
		"HS384": test384Key,
		"HS512": test512Key,
		"RS256": testRsaKey,
		"RS384": testRsaKey,
		"RS512": testRsaKey,
		"ES256": testEcdsa256Key,
		"ES384": testEcdsa384Key,
		"ES512": testEcdsa521Key,
		"PS256": testRsaKey,
		"PS384": testRsaKey,
		"PS512": testRsaKey,
	}
	testVeriKeys = map[string]interface{}{
		"none":  nil,
		"HS256": test256Key,
		"HS384": test384Key,
		"HS512": test512Key,
		"RS256": &testRsaKey.PublicKey,
		"RS384": &testRsaKey.PublicKey,
		"RS512": &testRsaKey.PublicKey,
		"ES256": &testEcdsa256Key.PublicKey,
		"ES384": &testEcdsa384Key.PublicKey,
		"ES512": &testEcdsa521Key.PublicKey,
		"PS256": &testRsaKey.PublicKey,
		"PS384": &testRsaKey.PublicKey,
		"PS512": &testRsaKey.PublicKey,
	}

	testEncKeys = map[string]interface{}{
		"RSA1_5":       &testRsaKey.PublicKey,
		"RSA-OAEP":     &testRsaKey.PublicKey,
		"RSA-OAEP-256": &testRsaKey.PublicKey,
		"A128KW":       test128Key,
		"A192KW":       test192Key,
		"A256KW":       test256Key,
		// "ECDH-ES":            &testEcdsa256Key.PublicKey,
		// "ECDH-ES+A128KW":     &testEcdsa256Key.PublicKey,
		// "ECDH-ES+A192KW":     &testEcdsa384Key.PublicKey,
		// "ECDH-ES+A256KW":     &testEcdsa521Key.PublicKey,
		"A128GCMKW": test128Key,
		"A192GCMKW": test192Key,
		"A256GCMKW": test256Key,
		// "PBES2-HS256+A128KW": test128Key,
		// "PBES2-HS384+A192KW": test192Key,
		// "PBES2-HS512+A256KW": test256Key,
	}
	testDecKeys = map[string]interface{}{
		"RSA1_5":       testRsaKey,
		"RSA-OAEP":     testRsaKey,
		"RSA-OAEP-256": testRsaKey,
		"A128KW":       test128Key,
		"A192KW":       test192Key,
		"A256KW":       test256Key,
		// "ECDH-ES":            testEcdsa256Key,
		// "ECDH-ES+A128KW":     testEcdsa256Key,
		// "ECDH-ES+A192KW":     testEcdsa384Key,
		// "ECDH-ES+A256KW":     testEcdsa521Key,
		"A128GCMKW": test128Key,
		"A192GCMKW": test192Key,
		"A256GCMKW": test256Key,
		// "PBES2-HS256+A128KW": test128Key,
		// "PBES2-HS384+A192KW": test192Key,
		// "PBES2-HS512+A256KW": test256Key,
	}
}

func TestRsaPublickKey(t *testing.T) {
	// JWK Appendix A.1 より。
	m := map[string]interface{}{
		"kty": "RSA",
		"n":   "0vx7agoebGcQSuuPiLJXZptN9nndrQmbXEps2aiAFbWhM78LhWx4cbbfAAtVT86zwu1RK7aPFFxuhDR1L6tSoc_BJECPebWKRXjBZCiFV4n3oknjhMstn64tZ_2W-5JsGY4Hc5n9yBXArwl93lqt7_RN5w6Cf0h4QyQ5v-65YGjQR0_FDW2QvzqY368QQMicAtaSqzs8KJZgnYb9c7d0zgdAZHzu6qMQvRL5hajrn1n91CbOpbISD08qNLyrdkt-bFTWhAI4vMQFh6WeZu0fM4lFd2NcRwr3XPksINHaQ-G_xBniIqbw0Ls1jF44-csFCur-kEgU8awapJzKnqDKgw",
		"e":   "AQAB",
	}

	if key, err := KeyFromJwkMap(m); err != nil {
		t.Fatal(err)
	} else if buff := KeyToJwkMap(key, nil); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestEcdsaPublicKey(t *testing.T) {
	// JWK Appendix A.1 より。
	m := map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y":   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
	}

	if key, err := KeyFromJwkMap(m); err != nil {
		t.Fatal(err)
	} else if buff := KeyToJwkMap(key, nil); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestRsaPrivateKey(t *testing.T) {
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

	if key, err := KeyFromJwkMap(m); err != nil {
		t.Fatal(err)
	} else if rsaKey, ok := key.(*rsa.PrivateKey); !ok {
		t.Error("not rsa private key")
		t.Error(key)
	} else if buff := KeyToJwkMap(key, nil); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	} else if decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, rsaKey, encrypted); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(decrypted, plain) {
		t.Error(decrypted)
		t.Error(plain)
	}
}

func TestEcdsaPrivateKey(t *testing.T) {
	// JWK Appendix A.2 より。

	m := map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "MKBCTNIcKUSDii11ySs3526iDZ8AiTo7Tu6KPAqv7D4",
		"y":   "4Etl6SRW2YiLUrN5vfvVHuhp7x8PxltmWWlbbM4IFyM",
		"d":   "870MB6gfuTJ4HtUnUvYMyJpr5eUZNP4Bk43bVdj3eAE",
	}

	if key, err := KeyFromJwkMap(m); err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*ecdsa.PrivateKey); !ok {
		t.Error("not ecdsa private key")
		t.Error(key)
	} else if buff := KeyToJwkMap(key, nil); !reflect.DeepEqual(m, buff) {
		t.Error(m)
		t.Error(buff)
	}
}

func TestJwk(t *testing.T) {
	for _, key := range []interface{}{
		test128Key,
		test192Key,
		test256Key,
		test384Key,
		test512Key,
		testRsaKey,
		testEcdsa256Key,
		testEcdsa384Key,
		testEcdsa521Key,
		&testRsaKey.PublicKey,
		&testEcdsa256Key.PublicKey,
		&testEcdsa384Key.PublicKey,
		&testEcdsa521Key.PublicKey,
	} {
		if key2, err := KeyFromJwkMap(KeyToJwkMap(key, nil)); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(key2, key) {
			t.Error(key2)
			t.Error(key)
		}
	}
}
