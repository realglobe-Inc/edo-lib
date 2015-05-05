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
	"bytes"
	"github.com/realglobe-Inc/edo-lib/jwk"
	"testing"
)

func TestParseJwsNoneSample(t *testing.T) {
	// JWT 6.1, JWS Appendix A.5 より。
	data := []byte("eyJhbGciOiJub25lIn0" +
		".eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		".")

	if jt, err := Parse(data); err != nil {
		t.Fatal(err)
	} else if !jt.IsSigned() {
		t.Fatal("cannot detect sign")
	} else if jt.IsEncrypted() {
		t.Fatal("detect encryption")
	} else if err := jt.Verify(nil); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsHsSample(t *testing.T) {
	// JWS Appendix A.1 より。
	key, err := jwk.FromMap(map[string]interface{}{
		"kty": "oct",
		"k":   "AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9" +
		".eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		".dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk")

	if jt, err := Parse(data); err != nil {
		t.Error(err)
		t.Fatal(string(data))
	} else if !jt.IsSigned() {
		t.Fatal("cannot detect sign")
	} else if jt.IsEncrypted() {
		t.Fatal("detect encryption")
	} else if err := jt.Verify([]jwk.Key{key}); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsRsSample(t *testing.T) {
	// JWT 3.1, JWS Appendix A.2 より。
	key, err := jwk.FromMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("eyJhbGciOiJSUzI1NiJ9" +
		".eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		".cC4hiUPoj9Eetdgtv3hF80EGrhuB__dzERat0XF9g2VtQgr9PJbu3XOiZj5RZmh7AAuHIm4Bh-0Qc_lF5YKt_O8W2Fp5jujGbds9uJdbF9CUAr7t1dnZcAcQjbKBYNX4BAynRFdiuB--f_nZLgrnbyTyWzO75vRK5h6xBArLIARNPvkSjtQBMHlb1L07Qe7K0GarZRmB_eSN9383LcOLn6_dO--xi12jzDwusC-eOkHWEsqtFZESc6BfI7noOPqvhJ1phCnvWh6IeYI2w9QOYEUipUTI8np6LbgGY9Fs98rqVt5AXLIhWkWywlVmtVrBp0igcN_IoypGlUPQGe77Rw")

	if jt, err := Parse(data); err != nil {
		t.Fatal(err)
	} else if !jt.IsSigned() {
		t.Fatal("cannot detect sign")
	} else if jt.IsEncrypted() {
		t.Fatal("detect encryption")
	} else if err := jt.Verify([]jwk.Key{key}); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsEsSample(t *testing.T) {
	// JWS Appendix A.3 より。
	key, err := jwk.FromMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
		"y":   "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("eyJhbGciOiJFUzI1NiJ9" +
		".eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		".DtEhU3ljbEg8L38VWAfUAqOyKAM6-Xx-F4GawxaepmXFCgfTjDxw5djxLa8ISlSApmWQxfKTUJqPP3-Kg6NU1Q")

	if jt, err := Parse(data); err != nil {
		t.Fatal(err)
	} else if !jt.IsSigned() {
		t.Fatal("cannot detect sign")
	} else if jt.IsEncrypted() {
		t.Fatal("detect encryption")
	} else if err := jt.Verify([]jwk.Key{key}); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJweRsa15Sample(t *testing.T) {
	// JWT Appendix A.1 より。
	key, err := jwk.FromMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "sXchDaQebHnPiGvyDOAT4saGEUetSyo9MKLOoWFsueri23bOdgWp4Dy1WlUzewbgBHod5pcM9H95GQRV3JDXboIRROSBigeC5yjU1hGzHHyXss8UDprecbAYxknTcQkhslANGRUZmdTOQ5qTRsLAt6BTYuyvVRdhS8exSZEy_c4gs_7svlJJQ4H9_NxsiIoLwAEk7-Q3UXERGYw_75IDrGA84-lA_-Ct4eTlXHBIY2EaV7t7LjJaynVJCpkv4LKjTTAumiGUIuQhrNhZLuF_RJLqHpM2kgWFLU7-VTdL1VbC2tejvcI2BlMkEpk1BzBZI0KQB0GaDWFLN-aEAw3vRw",
		"e":   "AQAB",
		"d":   "VFCWOqXr8nvZNyaaJLXdnNPXZKRaWCjkU5Q2egQQpTBMwhprMzWzpR8Sxq1OPThh_J6MUD8Z35wky9b8eEO0pwNS8xlh1lOFRRBoNqDIKVOku0aZb-rynq8cxjDTLZQ6Fz7jSjR1Klop-YKaUHc9GsEofQqYruPhzSA-QgajZGPbE_0ZaVDJHfyd7UUBUKunFMScbflYAAOYJqVIVwaYR5zWEEceUjNnTNo_CVSj-VvXLO5VZfCUAVLgW4dpf1SrtZjSt34YLsRarSb127reG_DUwg9Ch-KyvjT1SkHgUWRVGcyly7uvVGRSDwsXypdrNinPA4jlhoNdizK2zF2CWQ",
		"p":   "9gY2w6I6S6L0juEKsbeDAwpd9WMfgqFoeA9vEyEUuk4kLwBKcoe1x4HG68ik918hdDSE9vDQSccA3xXHOAFOPJ8R9EeIAbTi1VwBYnbTp87X-xcPWlEPkrdoUKW60tgs1aNd_Nnc9LEVVPMS390zbFxt8TN_biaBgelNgbC95sM",
		"q":   "uKlCKvKv_ZJMVcdIs5vVSU_6cPtYI1ljWytExV_skstvRSNi9r66jdd9-yBhVfuG4shsp2j7rGnIio901RBeHo6TPKWVVykPu1iYhQXw1jIABfw-MVsN-3bQ76WLdt2SDxsHs7q7zPyUyHXmps7ycZ5c72wGkUwNOjYelmkiNS0",
		"dp":  "w0kZbV63cVRvVX6yk3C8cMxo2qCM4Y8nsq1lmMSYhG4EcL6FWbX5h9yuvngs4iLEFk6eALoUS4vIWEwcL4txw9LsWH_zKI-hwoReoP77cOdSL4AVcraHawlkpyd2TWjE5evgbhWtOxnZee3cXJBkAi64Ik6jZxbvk-RR3pEhnCs",
		"dq":  "o_8V14SezckO6CNLKs_btPdFiO9_kC1DsuUTd2LAfIIVeMZ7jn1Gus_Ff7B7IVx3p5KuBGOVF8L-qifLb6nQnLysgHDh132NDioZkhH7mI7hPG-PYE_odApKdnqECHWw0J-F0JWnUd6D2B_1TvF9mXA2Qx-iGYn8OVV1Bsmp6qU",
		"qi":  "eNho5yRBEBxhGBtQRww9QirZsB66TrfFReG_CcteI1aCneT0ELGhYlRlCtUkTRclIfuEPmNsNDPbLoLqqCVznFbvdB7x-Tl-m0l_eFTj2KiqwGqE9PZB9nNTwMVvH3VRRSLWACvPnSiwP8N5Usy-WRXS-V7TbpxIhvepTfE0NNo",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("eyJhbGciOiJSU0ExXzUiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0" +
		".QR1Owv2ug2WyPBnbQrRARTeEk9kDO2w8qDcjiHnSJflSdv1iNqhWXaKH4MqAkQtMoNfABIPJaZm0HaA415sv3aeuBWnD8J-Ui7Ah6cWafs3ZwwFKDFUUsWHSK-IPKxLGTkND09XyjORj_CHAgOPJ-Sd8ONQRnJvWn_hXV1BNMHzUjPyYwEsRhDhzjAD26imasOTsgruobpYGoQcXUwFDn7moXPRfDE8-NoQX7N7ZYMmpUDkR-Cx9obNGwJQ3nM52YCitxoQVPzjbl7WBuB7AohdBoZOdZ24WlN1lVIeh8v1K4krB8xgKvRU8kgFrEn_a1rZgN5TiysnmzTROF869lQ" +
		".AxY8DCtDaGlsbGljb3RoZQ" +
		".MKOle7UQrG6nSxTLX6Mqwt0orbHvAKeWnDYvpIAeZ72deHxz3roJDXQyhxx0wKaMHDjUEOKIwrtkHthpqEanSBNYHZgmNOV7sln1Eu9g3J8" +
		".fiK51VwhsxJ-siBMR-YFiA")

	if jt, err := Parse(data); err != nil {
		t.Fatal(err)
	} else if jt.IsSigned() {
		t.Fatal("detect sign")
	} else if !jt.IsEncrypted() {
		t.Fatal("cannot detect encryption")
	} else if err := jt.Decrypt([]jwk.Key{key}); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJweNestingJwsSample(t *testing.T) {
	// JWT Appendix A.2 より。
	veriKey, err := jwk.FromMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	})
	if err != nil {
		t.Fatal(err)
	}
	decKey, err := jwk.FromMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "sXchDaQebHnPiGvyDOAT4saGEUetSyo9MKLOoWFsueri23bOdgWp4Dy1WlUzewbgBHod5pcM9H95GQRV3JDXboIRROSBigeC5yjU1hGzHHyXss8UDprecbAYxknTcQkhslANGRUZmdTOQ5qTRsLAt6BTYuyvVRdhS8exSZEy_c4gs_7svlJJQ4H9_NxsiIoLwAEk7-Q3UXERGYw_75IDrGA84-lA_-Ct4eTlXHBIY2EaV7t7LjJaynVJCpkv4LKjTTAumiGUIuQhrNhZLuF_RJLqHpM2kgWFLU7-VTdL1VbC2tejvcI2BlMkEpk1BzBZI0KQB0GaDWFLN-aEAw3vRw",
		"e":   "AQAB",
		"d":   "VFCWOqXr8nvZNyaaJLXdnNPXZKRaWCjkU5Q2egQQpTBMwhprMzWzpR8Sxq1OPThh_J6MUD8Z35wky9b8eEO0pwNS8xlh1lOFRRBoNqDIKVOku0aZb-rynq8cxjDTLZQ6Fz7jSjR1Klop-YKaUHc9GsEofQqYruPhzSA-QgajZGPbE_0ZaVDJHfyd7UUBUKunFMScbflYAAOYJqVIVwaYR5zWEEceUjNnTNo_CVSj-VvXLO5VZfCUAVLgW4dpf1SrtZjSt34YLsRarSb127reG_DUwg9Ch-KyvjT1SkHgUWRVGcyly7uvVGRSDwsXypdrNinPA4jlhoNdizK2zF2CWQ",
		"p":   "9gY2w6I6S6L0juEKsbeDAwpd9WMfgqFoeA9vEyEUuk4kLwBKcoe1x4HG68ik918hdDSE9vDQSccA3xXHOAFOPJ8R9EeIAbTi1VwBYnbTp87X-xcPWlEPkrdoUKW60tgs1aNd_Nnc9LEVVPMS390zbFxt8TN_biaBgelNgbC95sM",
		"q":   "uKlCKvKv_ZJMVcdIs5vVSU_6cPtYI1ljWytExV_skstvRSNi9r66jdd9-yBhVfuG4shsp2j7rGnIio901RBeHo6TPKWVVykPu1iYhQXw1jIABfw-MVsN-3bQ76WLdt2SDxsHs7q7zPyUyHXmps7ycZ5c72wGkUwNOjYelmkiNS0",
		"dp":  "w0kZbV63cVRvVX6yk3C8cMxo2qCM4Y8nsq1lmMSYhG4EcL6FWbX5h9yuvngs4iLEFk6eALoUS4vIWEwcL4txw9LsWH_zKI-hwoReoP77cOdSL4AVcraHawlkpyd2TWjE5evgbhWtOxnZee3cXJBkAi64Ik6jZxbvk-RR3pEhnCs",
		"dq":  "o_8V14SezckO6CNLKs_btPdFiO9_kC1DsuUTd2LAfIIVeMZ7jn1Gus_Ff7B7IVx3p5KuBGOVF8L-qifLb6nQnLysgHDh132NDioZkhH7mI7hPG-PYE_odApKdnqECHWw0J-F0JWnUd6D2B_1TvF9mXA2Qx-iGYn8OVV1Bsmp6qU",
		"qi":  "eNho5yRBEBxhGBtQRww9QirZsB66TrfFReG_CcteI1aCneT0ELGhYlRlCtUkTRclIfuEPmNsNDPbLoLqqCVznFbvdB7x-Tl-m0l_eFTj2KiqwGqE9PZB9nNTwMVvH3VRRSLWACvPnSiwP8N5Usy-WRXS-V7TbpxIhvepTfE0NNo",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("eyJhbGciOiJSU0ExXzUiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwiY3R5IjoiSldUIn0" +
		".g_hEwksO1Ax8Qn7HoN-BVeBoa8FXe0kpyk_XdcSmxvcM5_P296JXXtoHISr_DD_MqewaQSH4dZOQHoUgKLeFly-9RI11TG-_Ge1bZFazBPwKC5lJ6OLANLMd0QSL4fYEb9ERe-epKYE3xb2jfY1AltHqBO-PM6j23Guj2yDKnFv6WO72tteVzm_2n17SBFvhDuR9a2nHTE67pe0XGBUS_TK7ecA-iVq5COeVdJR4U4VZGGlxRGPLRHvolVLEHx6DYyLpw30Ay9R6d68YCLi9FYTq3hIXPK_-dmPlOUlKvPr1GgJzRoeC9G5qCvdcHWsqJGTO_z3Wfo5zsqwkxruxwA" +
		".UmVkbW9uZCBXQSA5ODA1Mg" +
		".VwHERHPvCNcHHpTjkoigx3_ExK0Qc71RMEParpatm0X_qpg-w8kozSjfNIPPXiTBBLXR65CIPkFqz4l1Ae9w_uowKiwyi9acgVztAi-pSL8GQSXnaamh9kX1mdh3M_TT-FZGQFQsFhu0Z72gJKGdfGE-OE7hS1zuBD5oEUfk0Dmb0VzWEzpxxiSSBbBAzP10l56pPfAtrjEYw-7ygeMkwBl6Z_mLS6w6xUgKlvW6ULmkV-uLC4FUiyKECK4e3WZYKw1bpgIqGYsw2v_grHjszJZ-_I5uM-9RA8ycX9KqPRp9gc6pXmoU_-27ATs9XCvrZXUtK2902AUzqpeEUJYjWWxSNsS-r1TJ1I-FMJ4XyAiGrfmo9hQPcNBYxPz3GQb28Y5CLSQfNgKSGt0A4isp1hBUXBHAndgtcslt7ZoQJaKe_nNJgNliWtWpJ_ebuOpEl8jdhehdccnRMIwAmU1n7SPkmhIl1HlSOpvcvDfhUN5wuqU955vOBvfkBOh5A11UzBuo2WlgZ6hYi9-e3w29bR0C2-pp3jbqxEDw3iWaf2dc5b-LnR0FEYXvI_tYk5rd_J9N0mg0tQ6RbpxNEMNoA9QWk5lgdPvbh9BaO195abQ" +
		".AVO9iT5AV4CzvDJCdhSFlQ")

	if jt, err := Parse(data); err != nil {
		t.Fatal(err)
	} else if jt.IsSigned() {
		t.Fatal("detect sign")
	} else if !jt.IsEncrypted() {
		t.Fatal("cannot detect encryption")
	} else if err := jt.Decrypt([]jwk.Key{decKey}); err != nil {
		t.Fatal(err)
	} else if jt2, err := Parse(jt.RawBody()); err != nil {
		t.Fatal(err)
	} else if !jt2.IsSigned() {
		t.Fatal("cannot detect sign")
	} else if jt2.IsEncrypted() {
		t.Fatal("detect encryption")
	} else if err := jt2.Verify([]jwk.Key{veriKey}); err != nil {
		t.Fatal(err)
	} else if jt2.Claim("iss") != "joe" {
		t.Fatal(jt2.Claim("iss"))
	} else if jt2.Claim("exp") != 1300819380.0 {
		t.Fatal(jt2.Claim("exp"))
	} else if jt2.Claim("http://example.com/is_root") != true {
		t.Fatal(jt2.Claim("http://example.com/is_root"))
	}
}

func TestJwt(t *testing.T) {
	rawHead := []byte(`{"alg":"none"}`)
	rawBody := []byte(`{"iss":"joe",
 "exp":1300819380,
 "http://example.com/is_root":true}`)

	jt := New()
	jt.SetRawHeader(rawHead)
	jt.SetRawBody(rawBody)

	if raw := jt.RawHeader(); !bytes.Equal(raw, rawHead) {
		t.Error(string(raw))
		t.Fatal(string(rawHead))
	} else if raw := jt.RawBody(); !bytes.Equal(raw, rawBody) {
		t.Error(string(raw))
		t.Fatal(string(rawBody))
	} else if jt.Header("alg") != "none" {
		t.Fatal(jt.Claim("alg"))
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}

	jt2 := New()
	jt2.SetHeader("alg", "none")
	jt2.SetClaim("iss", "joe")
	jt2.SetClaim("exp", 1300819380.0)
	jt2.SetClaim("http://example.com/is_root", true)
	if jt2.Header("alg") != "none" {
		t.Fatal(jt2.Claim("alg"))
	} else if jt2.Claim("iss") != "joe" {
		t.Fatal(jt2.Claim("iss"))
	} else if jt2.Claim("exp") != 1300819380.0 {
		t.Fatal(jt2.Claim("exp"))
	} else if jt2.Claim("http://example.com/is_root") != true {
		t.Fatal(jt2.Claim("http://example.com/is_root"))
	}
}

func TestJwtRemove(t *testing.T) {
	jt := New()
	jt.SetHeader("alg", "none")
	jt.SetHeader("test", "test")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("test", "test")

	if jt.Header("alg") != "none" {
		t.Fatal(jt.Header("alg"))
	} else if jt.Header("test") != "test" {
		t.Fatal(jt.Header("test"))
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("test") != "test" {
		t.Fatal(jt.Claim("test"))
	}

	jt.SetHeader("test", nil)
	jt.SetClaim("test", nil)

	if jt.Header("alg") != "none" {
		t.Fatal(jt.Header("alg"))
	} else if jt.Header("test") != nil {
		t.Fatal(jt.Header("test"))
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("test") != nil {
		t.Fatal(jt.Claim("test"))
	}
}

func TestJwtEncode(t *testing.T) {
	jt := New()
	jt.SetHeader("alg", "none")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380.0)
	jt.SetClaim("http://example.com/is_root", true)

	buff, err := jt.Encode()
	if err != nil {
		t.Fatal(err)
	}

	jt, err = Parse(buff)
	if err != nil {
		t.Error(err)
		t.Fatal(string(buff))
	}

	if jt.Header("alg") != "none" {
		t.Fatal(jt.Header("alg"))
	} else if jt.Claim("iss") != "joe" {
		t.Fatal(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Fatal(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Fatal(jt.Claim("http://example.com/is_root"))
	}
}

func TestJws(t *testing.T) {
	for _, key := range test_sigKeys {
		jt := New()
		jt.SetHeader("alg", key.Id())
		jt.SetHeader("kid", key.Id())
		jt.SetClaim("iss", "joe")
		jt.SetClaim("exp", 1300819380.0)
		jt.SetClaim("http://example.com/is_root", true)

		if err := jt.Sign(test_sigKeys); err != nil {
			t.Fatal(err)
		} else if buff, err := jt.Encode(); err != nil {
			t.Fatal(err)
		} else if jt2, err := Parse(buff); err != nil {
			t.Fatal(err)
		} else if !jt2.IsSigned() {
			t.Fatal("not signed")
		} else if jt2.IsEncrypted() {
			t.Fatal("encrypted")
		} else if err := jt2.Verify(test_veriKeys); err != nil {
			t.Fatal(err)
		} else if jt.Claim("iss") != "joe" {
			t.Fatal(jt.Claim("iss"))
		} else if jt.Claim("exp") != 1300819380.0 {
			t.Fatal(jt.Claim("exp"))
		} else if jt.Claim("http://example.com/is_root") != true {
			t.Fatal(jt.Claim("http://example.com/is_root"))
		}
	}
}

func TestJwe(t *testing.T) {
	for _, key := range test_encKeys {
		for _, enc := range []string{
			"A128CBC-HS256",
			"A192CBC-HS384",
			"A256CBC-HS512",
			"A128GCM",
			"A192GCM",
			"A256GCM",
		} {
			jt := New()
			jt.SetHeader("alg", key.Id())
			jt.SetHeader("kid", key.Id())
			jt.SetHeader("enc", enc)
			jt.SetClaim("iss", "joe")
			jt.SetClaim("exp", 1300819380.0)
			jt.SetClaim("http://example.com/is_root", true)

			if err := jt.Encrypt(test_encKeys); err != nil {
				t.Fatal(err)
			} else if buff, err := jt.Encode(); err != nil {
				t.Fatal(err)
			} else if jt2, err := Parse(buff); err != nil {
				t.Fatal(err)
			} else if jt2.IsSigned() {
				t.Fatal("signed")
			} else if !jt2.IsEncrypted() {
				t.Fatal("not encrypted")
			} else if err := jt2.Decrypt(test_decKeys); err != nil {
				t.Error(err)
				t.Fatal(string(jt.headPart))
			} else if jt2.Claim("iss") != "joe" {
				t.Fatal(jt2.Claim("iss"))
			} else if jt2.Claim("exp") != 1300819380.0 {
				t.Fatal(jt2.Claim("exp"))
			} else if jt2.Claim("http://example.com/is_root") != true {
				t.Fatal(jt2.Claim("http://example.com/is_root"))
			}
		}
	}
}

func TestJweDir(t *testing.T) {
	type param struct {
		enc string
		key jwk.Key
	}
	for _, p := range []param{
		{"A128CBC-HS256", test_256Key},
		{"A192CBC-HS384", test_384Key},
		{"A256CBC-HS512", test_512Key},
		{"A128GCM", test_128Key},
		{"A192GCM", test_192Key},
		{"A256GCM", test_256Key},
	} {
		jt := New()
		jt.SetHeader("alg", "dir")
		jt.SetHeader("enc", p.enc)
		jt.SetClaim("iss", "joe")
		jt.SetClaim("exp", 1300819380.0)
		jt.SetClaim("http://example.com/is_root", true)

		keys := []jwk.Key{p.key}
		if err := jt.Encrypt(keys); err != nil {
			t.Fatal(err)
		} else if buff, err := jt.Encode(); err != nil {
			t.Fatal(err)
		} else if jt2, err := Parse(buff); err != nil {
			t.Fatal(err)
		} else if jt2.IsSigned() {
			t.Fatal("signed")
		} else if !jt2.IsEncrypted() {
			t.Fatal("not encrypted")
		} else if err := jt2.Decrypt(keys); err != nil {
			t.Fatal(err)
		} else if jt2.Claim("iss") != "joe" {
			t.Fatal(jt2.Claim("iss"))
		} else if jt2.Claim("exp") != 1300819380.0 {
			t.Fatal(jt2.Claim("exp"))
		} else if jt2.Claim("http://example.com/is_root") != true {
			t.Fatal(jt2.Claim("http://example.com/is_root"))
		}
	}

}

func TestJweNestingJws(t *testing.T) {
	nested := New()
	nested.SetHeader("alg", "ES512")
	nested.SetHeader("kid", "ES512")
	nested.SetClaim("iss", "joe")
	nested.SetClaim("exp", 1300819380.0)
	nested.SetClaim("http://example.com/is_root", true)

	if err := nested.Sign(test_sigKeys); err != nil {
		t.Fatal(err)
	}
	buff, err := nested.Encode()
	if err != nil {
		t.Fatal(err)
	}

	jt := New()
	jt.SetHeader("alg", "RSA-OAEP-256")
	jt.SetHeader("kid", "RSA-OAEP-256")
	jt.SetHeader("enc", "A256CBC-HS512")
	jt.SetRawBody(buff)

	if err := jt.Encrypt(test_encKeys); err != nil {
		t.Fatal(err)
	} else if buff2, err := jt.Encode(); err != nil {
		t.Fatal(err)
	} else if jt2, err := Parse(buff2); err != nil {
		t.Fatal(err)
	} else if jt2.IsSigned() {
		t.Fatal("signed")
	} else if !jt2.IsEncrypted() {
		t.Fatal("not encrypted")
	} else if err := jt2.Decrypt(test_decKeys); err != nil {
		t.Fatal(err)
	} else if jt3, err := Parse(jt2.RawBody()); err != nil {
		t.Fatal(err)
	} else if !jt3.IsSigned() {
		t.Fatal("not signed")
	} else if jt3.IsEncrypted() {
		t.Fatal("encrypted")
	} else if err := jt3.Verify(test_veriKeys); err != nil {
		t.Fatal(err)
	} else if jt3.Claim("iss") != "joe" {
		t.Fatal(jt3.Claim("iss"))
	} else if jt3.Claim("exp") != 1300819380.0 {
		t.Fatal(jt3.Claim("exp"))
	} else if jt3.Claim("http://example.com/is_root") != true {
		t.Fatal(jt3.Claim("http://example.com/is_root"))
	}
}

func TestJweZip(t *testing.T) {
	jt := New()
	jt.SetHeader("alg", "RSA-OAEP-256")
	jt.SetHeader("kid", "RSA-OAEP-256")
	jt.SetHeader("enc", "A256CBC-HS512")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380.0)
	jt.SetClaim("jti", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	jt.SetClaim("http://example.com/is_root", true)

	if err := jt.Encrypt(test_encKeys); err != nil {
		t.Fatal(err)
	}
	buff0, err := jt.Encode()
	if err != nil {
		t.Fatal(err)
	}

	for zip := range []string{"DEF"} {
		jt.SetHeader("zip", zip)
		if err := jt.Encrypt(test_encKeys); err != nil {
			t.Fatal(err)
		} else if buff, err := jt.Encode(); err != nil {
			t.Fatal(err)
		} else if len(buff0) >= len(buff) {
			t.Error(string(buff0))
			t.Fatal(string(buff))
		} else if jt2, err := Parse(buff); err != nil {
			t.Fatal(err)
		} else if jt2.IsSigned() {
			t.Fatal("signed")
		} else if !jt2.IsEncrypted() {
			t.Fatal("not encrypted")
		} else if err := jt2.Decrypt(test_decKeys); err != nil {
			t.Fatal(err)
		} else if jt2.Claim("iss") != "joe" {
			t.Fatal(jt2.Claim("iss"))
		} else if jt2.Claim("exp") != 1300819380.0 {
			t.Fatal(jt2.Claim("exp"))
		} else if jt2.Claim("http://example.com/is_root") != true {
			t.Fatal(jt2.Claim("http://example.com/is_root"))
		}
	}
}
