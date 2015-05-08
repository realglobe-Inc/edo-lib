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
	"github.com/realglobe-Inc/edo-lib/jwk"
)

var test_128Key, test_192Key, test_256Key, test_384Key, test_512Key jwk.Key
var test_rsaKey jwk.Key
var test_ec256Key, test_ec384Key, test_ec521Key jwk.Key
var test_ecKey jwk.Key

var test_sigKeys, test_veriKeys []jwk.Key
var test_encKeys, test_decKeys []jwk.Key

func init() {
	buff := []byte{}
	for ; len(buff) < 64; buff = append(buff, byte(len(buff))) {
	}
	test_128Key = jwk.New(buff[:16], nil)
	test_192Key = jwk.New(buff[:24], nil)
	test_256Key = jwk.New(buff[:32], nil)
	test_384Key = jwk.New(buff[:48], nil)
	test_512Key = jwk.New(buff, nil)

	var err error
	test_rsaKey, err = jwk.FromMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "5cadP6Vvv6ABglXpSeXYxPB321gtSwmjccsHr2-YKmBm22KWF2A1b68LJ3mA8eG5NPSRL6macCMttxsoAKwaCxOxn-6dNOKXNLQ1S0WsE4yY2QLoi9Cj_sY8yfdk_wb0ZM5kyE99GjFFLDvnh-RjHIf2cbXPyPfbeLigeeon7jsxOw",
		"e":   "AQAB",
		"d":   "gOV1-Oo5UenUbuT6xXWmsHOlCOriHaH-iis22HdliQAjMxaO0_Yog8pSG4bRit7xIn-_olkmRZm2X21gd2AUC_mkE7Nytw5t_pioMzupEVVGApIFuc2_ryf5VPSznx3zk5FY6XCgUf6BnJ188WRUv3CnnNuAmEJtP6MhWmoKlPMpgQ",
		"p":   "9556qFgzilKEEhQ41fVzvLm5vKpiCc0IABG1CDQ_VTr4KGoOcqSHx6__yqYFQlzgizkG-zVxBQSs-6GZ3eA-t4s",
		"q":   "7Y2H2tRgIm9UjN0OlszOBcOXqPicE5KlseuNCIJZo1SyW30h-N2ssjCeiSDPrqm5QGZ637EAmhvNsPNOxzwLIxE",
		"dp":  "sa3EMdvoT87Z-ecMyWpw-_EA-AICiynWHcaW8iYbc9r2inlfmJ-61mzRzOXITFA8x2nKOqOkT4eFYKIauHzaQ_U",
		"dq":  "EhoQ2ioI0VbueHV34SHmKSZIbkXTjuJD4hTzAEz-i6Wuma4lYpNxz3pI-mYXrVWdmjy07ErOou-vcuZ3gFMg_iE",
		"qi":  "DbisQAteFbdCaNy6TyNy5UgZjdPba1bhKI3iIXalno_5HRrK4tUzu9VHdYVj5-iscIw5za9cPMLFr3zQvWa-gzA",
	})
	if err != nil {
		panic(err)
	}

	test_ec256Key, err = jwk.FromMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "lpHYO1qpjU95B2sThPR2-1jv44axgaEDkQtcKNE-oZs",
		"y":   "soy5O11SFFFeYdhQVodXlYPIpeo0pCS69IxiVPPf0Tk",
		"d":   "3BhkCluOkm8d8gvaPD5FDG2zeEw2JKf3D5LwN-mYmsw",
	})
	if err != nil {
		panic(err)
	}

	test_ec384Key, err = jwk.FromMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-384",
		"x":   "HlrMhzZww_AkmHV-2gDR5n7t75673UClnC7V2GewWva_sg-4GSUguFalVgwnK0tQ",
		"y":   "fxS48Fy50SZFZ-RAQRWUZXZgRSWwiKVkqPTd6gypfpQNkXSwE69BXYIAQcfaLcf2",
		"d":   "Gp-7eC0G7PjGzKoiAmTQ1iLsLU3AEy3h-bKFWSZOanXqSWI6wqJVPEUsatNYBJoG",
	})
	if err != nil {
		panic(err)
	}

	test_ec521Key, err = jwk.FromMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-521",
		"x":   "AXA9Y2pY1g_Cs4W6Nto7ebjDKOsaRHxo6EYRWjk1XHZaA7HnkeHg13x24OWHelqdZiuo7J1VbRKJ4ohZPjKX-AL7",
		"y":   "AeC5zdcUHLDdhQAvdsnnwD8rgNWjMdlWsqXZxv7Ar7ly5xmZGxEDtcJuhfhn8R9PXeScPH2soF3dFYPCuDkF4Gns",
		"d":   "AK9ejP0HuUt7ojjI9p20986DGqG-5jc9UWMtnMqxNvIvBScTJflS2lE-6CRsJZKk6ChWI6U4ahXDH0cCCFWTAvI_",
	})
	if err != nil {
		panic(err)
	}

	test_ecKey = test_ec256Key

	test_sigKeys = []jwk.Key{
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "HS256"}),
		jwk.New(test_384Key.Common(), map[string]interface{}{"kid": "HS384"}),
		jwk.New(test_512Key.Common(), map[string]interface{}{"kid": "HS512"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RS256"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RS384"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RS512"}),
		jwk.New(test_ec256Key.Private(), map[string]interface{}{"kid": "ES256"}),
		jwk.New(test_ec384Key.Private(), map[string]interface{}{"kid": "ES384"}),
		jwk.New(test_ec521Key.Private(), map[string]interface{}{"kid": "ES512"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "PS256"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "PS384"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "PS512"}),
	}
	test_veriKeys = []jwk.Key{
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "HS256"}),
		jwk.New(test_384Key.Common(), map[string]interface{}{"kid": "HS384"}),
		jwk.New(test_512Key.Common(), map[string]interface{}{"kid": "HS512"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RS256"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RS384"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RS512"}),
		jwk.New(test_ec256Key.Public(), map[string]interface{}{"kid": "ES256"}),
		jwk.New(test_ec384Key.Public(), map[string]interface{}{"kid": "ES384"}),
		jwk.New(test_ec521Key.Public(), map[string]interface{}{"kid": "ES512"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "PS256"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "PS384"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "PS512"}),
	}

	test_encKeys = []jwk.Key{
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RSA1_5"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RSA-OAEP"}),
		jwk.New(test_rsaKey.Public(), map[string]interface{}{"kid": "RSA-OAEP-256"}),
		jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "A128KW"}),
		jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "A192KW"}),
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "A256KW"}),
		// jwk.New(test_ec256Key.Public(), map[string]interface{}{"kid": "ECDH-ES"}),
		// jwk.New(test_ec256Key.Public(), map[string]interface{}{"kid": "ECDH-ES+A128KW"}),
		// jwk.New(test_ec384Key.Public(), map[string]interface{}{"kid": "ECDH-ES+A192KW"}),
		// jwk.New(test_ec521Key.Public(), map[string]interface{}{"kid": "ECDH-ES+A256KW"}),
		jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "A128GCMKW"}),
		jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "A192GCMKW"}),
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "A256GCMKW"}),
		// jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "PBES2-HS256+A128KW"}),
		// jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "PBES2-HS384+A192KW"}),
		// jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "PBES2-HS512+A256KW"}),
	}
	test_decKeys = []jwk.Key{
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RSA1_5"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RSA-OAEP"}),
		jwk.New(test_rsaKey.Private(), map[string]interface{}{"kid": "RSA-OAEP-256"}),
		jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "A128KW"}),
		jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "A192KW"}),
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "A256KW"}),
		// jwk.New(test_ec256Key.Private(), map[string]interface{}{"kid": "ECDH-ES"}),
		// jwk.New(test_ec256Key.Private(), map[string]interface{}{"kid": "ECDH-ES+A128KW"}),
		// jwk.New(test_ec384Key.Private(), map[string]interface{}{"kid": "ECDH-ES+A192KW"}),
		// jwk.New(test_ec521Key.Private(), map[string]interface{}{"kid": "ECDH-ES+A256KW"}),
		jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "A128GCMKW"}),
		jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "A192GCMKW"}),
		jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "A256GCMKW"}),
		// jwk.New(test_128Key.Common(), map[string]interface{}{"kid": "PBES2-HS256+A128KW"}),
		// jwk.New(test_192Key.Common(), map[string]interface{}{"kid": "PBES2-HS384+A192KW"}),
		// jwk.New(test_256Key.Common(), map[string]interface{}{"kid": "PBES2-HS512+A256KW"}),
	}
}
