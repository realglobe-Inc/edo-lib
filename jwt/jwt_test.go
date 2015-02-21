package jwt

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseJwsNoneSample(t *testing.T) {
	// JWT 6.1, JWS Appendix A.5 より。
	raw := "eyJhbGciOiJub25lIn0" +
		"." +
		"eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"." +
		""

	if jt, err := Parse(raw, nil, nil); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsHsSample(t *testing.T) {
	// JWS Appendix A.1 より。
	key, err := KeyFromJwkMap(map[string]interface{}{
		"kty": "oct",
		"k":   "AyM1SysPpbyDfgZld3umj1qzKObwVMkoqQ-EstJQLr_T-1qS0gZH75aKtMN3Yj0iPS4hcgUuTwjAzZr1Z9CAow",
	})
	if err != nil {
		t.Fatal(err)
	}
	raw := "eyJ0eXAiOiJKV1QiLA0KICJhbGciOiJIUzI1NiJ9" +
		"." +
		"eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"." +
		"dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"

	if jt, err := Parse(raw, map[string]interface{}{"": key}, nil); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsRsSample(t *testing.T) {
	// JWT 3.1, JWS Appendix A.2 より。
	key, err := KeyFromJwkMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	})
	if err != nil {
		t.Fatal(err)
	}
	raw := "eyJhbGciOiJSUzI1NiJ9" +
		"." +
		"eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"." +
		"cC4hiUPoj9Eetdgtv3hF80EGrhuB__dzERat0XF9g2VtQgr9PJbu3XOiZj5RZmh7AAuHIm4Bh-0Qc_lF5YKt_O8W2Fp5jujGbds9uJdbF9CUAr7t1dnZcAcQjbKBYNX4BAynRFdiuB--f_nZLgrnbyTyWzO75vRK5h6xBArLIARNPvkSjtQBMHlb1L07Qe7K0GarZRmB_eSN9383LcOLn6_dO--xi12jzDwusC-eOkHWEsqtFZESc6BfI7noOPqvhJ1phCnvWh6IeYI2w9QOYEUipUTI8np6LbgGY9Fs98rqVt5AXLIhWkWywlVmtVrBp0igcN_IoypGlUPQGe77Rw"

	if jt, err := Parse(raw, map[string]interface{}{"": key}, nil); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJwsEsSample(t *testing.T) {
	// JWS Appendix A.3 より。
	key, err := KeyFromJwkMap(map[string]interface{}{
		"kty": "EC",
		"crv": "P-256",
		"x":   "f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
		"y":   "x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0",
	})
	if err != nil {
		t.Fatal(err)
	}
	raw := "eyJhbGciOiJFUzI1NiJ9" +
		"." +
		"eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"." +
		"DtEhU3ljbEg8L38VWAfUAqOyKAM6-Xx-F4GawxaepmXFCgfTjDxw5djxLa8ISlSApmWQxfKTUJqPP3-Kg6NU1Q"

	if jt, err := Parse(raw, map[string]interface{}{"": key}, nil); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJweRsa15Sample(t *testing.T) {
	// JWT Appendix A.1 より。
	key, err := KeyFromJwkMap(map[string]interface{}{
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
	raw := "eyJhbGciOiJSU0ExXzUiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2In0" +
		"." +
		"QR1Owv2ug2WyPBnbQrRARTeEk9kDO2w8qDcjiHnSJflSdv1iNqhWXaKH4MqAkQtMoNfABIPJaZm0HaA415sv3aeuBWnD8J-Ui7Ah6cWafs3ZwwFKDFUUsWHSK-IPKxLGTkND09XyjORj_CHAgOPJ-Sd8ONQRnJvWn_hXV1BNMHzUjPyYwEsRhDhzjAD26imasOTsgruobpYGoQcXUwFDn7moXPRfDE8-NoQX7N7ZYMmpUDkR-Cx9obNGwJQ3nM52YCitxoQVPzjbl7WBuB7AohdBoZOdZ24WlN1lVIeh8v1K4krB8xgKvRU8kgFrEn_a1rZgN5TiysnmzTROF869lQ" +
		"." +
		"AxY8DCtDaGlsbGljb3RoZQ" +
		"." +
		"MKOle7UQrG6nSxTLX6Mqwt0orbHvAKeWnDYvpIAeZ72deHxz3roJDXQyhxx0wKaMHDjUEOKIwrtkHthpqEanSBNYHZgmNOV7sln1Eu9g3J8" +
		"." +
		"fiK51VwhsxJ-siBMR-YFiA"

	if jt, err := Parse(raw, nil, map[string]interface{}{"": key}); err != nil {
		t.Fatal(err)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestParseJweNestingJwsSample(t *testing.T) {
	// JWT Appendix A.2 より。
	veriKey, err := KeyFromJwkMap(map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	})
	if err != nil {
		t.Fatal(err)
	}
	decKey, err := KeyFromJwkMap(map[string]interface{}{
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
	raw := "eyJhbGciOiJSU0ExXzUiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwiY3R5IjoiSldUIn0" +
		"." +
		"g_hEwksO1Ax8Qn7HoN-BVeBoa8FXe0kpyk_XdcSmxvcM5_P296JXXtoHISr_DD_MqewaQSH4dZOQHoUgKLeFly-9RI11TG-_Ge1bZFazBPwKC5lJ6OLANLMd0QSL4fYEb9ERe-epKYE3xb2jfY1AltHqBO-PM6j23Guj2yDKnFv6WO72tteVzm_2n17SBFvhDuR9a2nHTE67pe0XGBUS_TK7ecA-iVq5COeVdJR4U4VZGGlxRGPLRHvolVLEHx6DYyLpw30Ay9R6d68YCLi9FYTq3hIXPK_-dmPlOUlKvPr1GgJzRoeC9G5qCvdcHWsqJGTO_z3Wfo5zsqwkxruxwA" +
		"." +
		"UmVkbW9uZCBXQSA5ODA1Mg" +
		"." +
		"VwHERHPvCNcHHpTjkoigx3_ExK0Qc71RMEParpatm0X_qpg-w8kozSjfNIPPXiTBBLXR65CIPkFqz4l1Ae9w_uowKiwyi9acgVztAi-pSL8GQSXnaamh9kX1mdh3M_TT-FZGQFQsFhu0Z72gJKGdfGE-OE7hS1zuBD5oEUfk0Dmb0VzWEzpxxiSSBbBAzP10l56pPfAtrjEYw-7ygeMkwBl6Z_mLS6w6xUgKlvW6ULmkV-uLC4FUiyKECK4e3WZYKw1bpgIqGYsw2v_grHjszJZ-_I5uM-9RA8ycX9KqPRp9gc6pXmoU_-27ATs9XCvrZXUtK2902AUzqpeEUJYjWWxSNsS-r1TJ1I-FMJ4XyAiGrfmo9hQPcNBYxPz3GQb28Y5CLSQfNgKSGt0A4isp1hBUXBHAndgtcslt7ZoQJaKe_nNJgNliWtWpJ_ebuOpEl8jdhehdccnRMIwAmU1n7SPkmhIl1HlSOpvcvDfhUN5wuqU955vOBvfkBOh5A11UzBuo2WlgZ6hYi9-e3w29bR0C2-pp3jbqxEDw3iWaf2dc5b-LnR0FEYXvI_tYk5rd_J9N0mg0tQ6RbpxNEMNoA9QWk5lgdPvbh9BaO195abQ" +
		"." +
		"AVO9iT5AV4CzvDJCdhSFlQ"

	if jt, err := Parse(raw, map[string]interface{}{"": veriKey}, map[string]interface{}{"": decKey}); err != nil {
		t.Fatal(err)
	} else if nested := jt.Nested(); nested == nil {
		t.Error("no nested JWT")
	} else if nested.Claim("iss") != "joe" {
		t.Error(nested.Claim("iss"))
	} else if nested.Claim("exp") != 1300819380.0 {
		t.Error(nested.Claim("exp"))
	} else if nested.Claim("http://example.com/is_root") != true {
		t.Error(nested.Claim("http://example.com/is_root"))
	}
}

func TestJwt(t *testing.T) {
	jt := New()
	jt.SetHeader("alg", "none")
	jt.SetHeader("test", "test")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("test", "test")

	if jt.Header("alg") != "none" {
		t.Error(jt.Header("alg"))
	} else if jt.Header("test") != "test" {
		t.Error(jt.Header("test"))
	} else if heads := jt.HeaderNames(); !heads["alg"] || !heads["test"] {
		t.Error(heads)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("test") != "test" {
		t.Error(jt.Claim("test"))
	} else if clms := jt.ClaimNames(); !clms["iss"] || !clms["test"] {
		t.Error(clms)
	}

	jt.SetHeader("test", "")
	jt.SetClaim("test", nil)

	if jt.Header("alg") != "none" {
		t.Error(jt.Header("alg"))
	} else if jt.Header("test") != nil {
		t.Error(jt.Header("test"))
	} else if heads := jt.HeaderNames(); !heads["alg"] || heads["test"] {
		t.Error(heads)
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("test") != nil {
		t.Error(jt.Claim("test"))
	} else if clms := jt.ClaimNames(); !clms["iss"] || clms["test"] {
		t.Error(clms)
	}
}

func TestJwtEncode(t *testing.T) {
	jt := New()
	jt.SetHeader("alg", "none")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380)
	jt.SetClaim("http://example.com/is_root", true)

	buff, err := jt.Encode(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	jt, err = Parse(string(buff), nil, nil)
	if err != nil {
		t.Fatal(err, string(buff))
	}

	if jt.Header("alg") != "none" {
		t.Error(jt.Header("alg"))
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380.0 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestJws(t *testing.T) {
	for alg := range testSigKeys {
		jt := New()
		jt.SetHeader("alg", alg)
		jt.SetHeader("kid", alg)
		jt.SetClaim("iss", "joe")
		jt.SetClaim("exp", 1300819380.0)
		jt.SetClaim("http://example.com/is_root", true)

		if encoded, err := jt.Encode(testSigKeys, nil); err != nil {
			t.Fatal(err)
		} else if jt2, err := Parse(string(encoded), testVeriKeys, nil); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(jt2, jt) {
			t.Error(fmt.Sprintf("%#v", jt2))
			t.Error(fmt.Sprintf("%#v", jt))
		}
	}
}

func TestJwe(t *testing.T) {
	for alg := range testEncKeys {
		for enc := range jweEncs {
			jt := New()
			jt.SetHeader("alg", alg)
			jt.SetHeader("kid", alg)
			jt.SetHeader("enc", enc)
			jt.SetClaim("iss", "joe")
			jt.SetClaim("exp", 1300819380.0)
			jt.SetClaim("http://example.com/is_root", true)

			if encoded, err := jt.Encode(nil, testEncKeys); err != nil {
				t.Error(alg + " " + enc)
				t.Fatal(err)
			} else if jt2, err := Parse(string(encoded), nil, testDecKeys); err != nil {
				t.Error(alg + " " + enc)
				t.Fatal(err)
			} else if !reflect.DeepEqual(jt2, jt) {
				t.Error(fmt.Sprintf("%#v", jt2))
				t.Error(fmt.Sprintf("%#v", jt))
			}
		}
	}
}

func TestJweDir(t *testing.T) {
	type param struct {
		enc string
		key []byte
	}
	for _, p := range []param{
		{"A128CBC-HS256", test256Key},
		{"A192CBC-HS384", test384Key},
		{"A256CBC-HS512", test512Key},
		{"A128GCM", test128Key},
		{"A192GCM", test192Key},
		{"A256GCM", test256Key},
	} {
		jt := New()
		jt.SetHeader("alg", "dir")
		jt.SetHeader("enc", p.enc)
		jt.SetClaim("iss", "joe")
		jt.SetClaim("exp", 1300819380.0)
		jt.SetClaim("http://example.com/is_root", true)

		keys := map[string]interface{}{"": p.key}
		if encoded, err := jt.Encode(nil, keys); err != nil {
			t.Error(p.enc)
			t.Fatal(err)
		} else if jt2, err := Parse(string(encoded), nil, keys); err != nil {
			t.Error(p.enc)
			t.Fatal(err)
		} else if !reflect.DeepEqual(jt2, jt) {
			t.Error(fmt.Sprintf("%#v", jt2))
			t.Error(fmt.Sprintf("%#v", jt))
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

	jt := New()
	jt.SetHeader("alg", "RSA-OAEP-256")
	jt.SetHeader("kid", "RSA-OAEP-256")
	jt.SetHeader("enc", "A256CBC-HS512")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380.0)
	jt.SetClaim("http://example.com/is_root", true)
	jt.Nest(nested)

	if encoded, err := jt.Encode(testSigKeys, testEncKeys); err != nil {
		t.Fatal(err)
	} else if jt2, err := Parse(string(encoded), testVeriKeys, testDecKeys); err != nil {
		t.Fatal(err)
	} else if !reflect.DeepEqual(jt2, jt) {
		t.Error(fmt.Sprintf("%#v", jt2))
		t.Error(fmt.Sprintf("%#v", jt))
	} else if !reflect.DeepEqual(jt2.Nested(), nested) {
		t.Error(fmt.Sprintf("%#v", jt2.Nested()))
		t.Error(fmt.Sprintf("%#v", nested))
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

	encoded, err := jt.Encode(testSigKeys, testEncKeys)
	if err != nil {
		t.Fatal(err)
	}

	for zip := range jweZips {
		jt.SetHeader("zip", zip)
		if encoded2, err := jt.Encode(testSigKeys, testEncKeys); err != nil {
			t.Fatal(err)
		} else if jt2, err := Parse(string(encoded2), testVeriKeys, testDecKeys); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(jt2, jt) {
			t.Error(fmt.Sprintf("%#v", jt2))
			t.Error(fmt.Sprintf("%#v", jt))
		} else if len(encoded2) >= len(encoded) {
			t.Error(string(encoded2))
			t.Error(string(encoded))
		}
	}
}
