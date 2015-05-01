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

package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/realglobe-Inc/edo-lib/jwk"
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

	if key, err := jwk.FromMap(map[string]interface{}{
		"d":   "gOV1-Oo5UenUbuT6xXWmsHOlCOriHaH-iis22HdliQAjMxaO0_Yog8pSG4bRit7xIn-_olkmRZm2X21gd2AUC_mkE7Nytw5t_pioMzupEVVGApIFuc2_ryf5VPSznx3zk5FY6XCgUf6BnJ188WRUv3CnnNuAmEJtP6MhWmoKlPMpgQ",
		"dp":  "sa3EMdvoT87Z-ecMyWpw-_EA-AICiynWHcaW8iYbc9r2inlfmJ-61mzRzOXITFA8x2nKOqOkT4eFYKIauHzaQ_U",
		"dq":  "EhoQ2ioI0VbueHV34SHmKSZIbkXTjuJD4hTzAEz-i6Wuma4lYpNxz3pI-mYXrVWdmjy07ErOou-vcuZ3gFMg_iE",
		"e":   "AQAB",
		"kty": "RSA",
		"n":   "5cadP6Vvv6ABglXpSeXYxPB321gtSwmjccsHr2-YKmBm22KWF2A1b68LJ3mA8eG5NPSRL6macCMttxsoAKwaCxOxn-6dNOKXNLQ1S0WsE4yY2QLoi9Cj_sY8yfdk_wb0ZM5kyE99GjFFLDvnh-RjHIf2cbXPyPfbeLigeeon7jsxOw",
		"p":   "9556qFgzilKEEhQ41fVzvLm5vKpiCc0IABG1CDQ_VTr4KGoOcqSHx6__yqYFQlzgizkG-zVxBQSs-6GZ3eA-t4s",
		"q":   "7Y2H2tRgIm9UjN0OlszOBcOXqPicE5KlseuNCIJZo1SyW30h-N2ssjCeiSDPrqm5QGZ637EAmhvNsPNOxzwLIxE",
		"qi":  "DbisQAteFbdCaNy6TyNy5UgZjdPba1bhKI3iIXalno_5HRrK4tUzu9VHdYVj5-iscIw5za9cPMLFr3zQvWa-gzA",
	}); err != nil {
		panic(err)
	} else {
		testRsaKey = key.Private().(*rsa.PrivateKey)
	}

	if key, err := jwk.FromMap(map[string]interface{}{
		"crv": "P-256",
		"d":   "3BhkCluOkm8d8gvaPD5FDG2zeEw2JKf3D5LwN-mYmsw",
		"kty": "EC",
		"x":   "lpHYO1qpjU95B2sThPR2-1jv44axgaEDkQtcKNE-oZs",
		"y":   "soy5O11SFFFeYdhQVodXlYPIpeo0pCS69IxiVPPf0Tk",
	}); err != nil {
		panic(err)
	} else {
		testEcdsa256Key = key.Private().(*ecdsa.PrivateKey)
	}

	if key, err := jwk.FromMap(map[string]interface{}{
		"crv": "P-384",
		"d":   "Gp-7eC0G7PjGzKoiAmTQ1iLsLU3AEy3h-bKFWSZOanXqSWI6wqJVPEUsatNYBJoG",
		"kty": "EC",
		"x":   "HlrMhzZww_AkmHV-2gDR5n7t75673UClnC7V2GewWva_sg-4GSUguFalVgwnK0tQ",
		"y":   "fxS48Fy50SZFZ-RAQRWUZXZgRSWwiKVkqPTd6gypfpQNkXSwE69BXYIAQcfaLcf2",
	}); err != nil {
		panic(err)
	} else {
		testEcdsa384Key = key.Private().(*ecdsa.PrivateKey)
	}

	if key, err := jwk.FromMap(map[string]interface{}{
		"crv": "P-521",
		"d":   "AK9ejP0HuUt7ojjI9p20986DGqG-5jc9UWMtnMqxNvIvBScTJflS2lE-6CRsJZKk6ChWI6U4ahXDH0cCCFWTAvI_",
		"kty": "EC",
		"x":   "AXA9Y2pY1g_Cs4W6Nto7ebjDKOsaRHxo6EYRWjk1XHZaA7HnkeHg13x24OWHelqdZiuo7J1VbRKJ4ohZPjKX-AL7",
		"y":   "AeC5zdcUHLDdhQAvdsnnwD8rgNWjMdlWsqXZxv7Ar7ly5xmZGxEDtcJuhfhn8R9PXeScPH2soF3dFYPCuDkF4Gns",
	}); err != nil {
		panic(err)
	} else {
		testEcdsa521Key = key.Private().(*ecdsa.PrivateKey)
	}

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
