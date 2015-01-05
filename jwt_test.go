package util

import (
	"testing"
)

func TestParseJwt(t *testing.T) {
	// JSON Web Token (JWT) より。
	var raw string = "eyJhbGciOiJub25lIn0" +
		"." + "eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"."

	jt, err := ParseJwt(raw)
	if err != nil {
		t.Fatal(err)
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

func TestJwt(t *testing.T) {
	jt := NewJwt()
	jt.SetHeader("alg", "none")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380)
	jt.SetClaim("http://example.com/is_root", true)

	if jt.Header("alg") != "none" {
		t.Error(jt.Header("alg"))
	} else if jt.Claim("iss") != "joe" {
		t.Error(jt.Claim("iss"))
	} else if jt.Claim("exp") != 1300819380 {
		t.Error(jt.Claim("exp"))
	} else if jt.Claim("http://example.com/is_root") != true {
		t.Error(jt.Claim("http://example.com/is_root"))
	}
}

func TestJwtEncode(t *testing.T) {
	jt := NewJwt()
	jt.SetHeader("alg", "none")
	jt.SetClaim("iss", "joe")
	jt.SetClaim("exp", 1300819380)
	jt.SetClaim("http://example.com/is_root", true)

	buff, err := jt.Encode()
	if err != nil {
		t.Fatal(err)
	}

	jt, err = ParseJwt(string(buff))
	if err != nil {
		t.Fatal(err)
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
