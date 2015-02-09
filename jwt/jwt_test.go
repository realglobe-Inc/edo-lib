package jwt

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseJwt(t *testing.T) {
	// JSON Web Token (JWT) より。
	var raw string = "eyJhbGciOiJub25lIn0" +
		"." + "eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"."

	jw, err := ParseJwt(raw)
	if err != nil {
		t.Fatal(err)
	}

	if jw.Header("alg") != "none" {
		t.Error(jw.Header("alg"))
	} else if jw.Claim("iss") != "joe" {
		t.Error(jw.Claim("iss"))
	} else if jw.Claim("exp") != 1300819380.0 {
		t.Error(jw.Claim("exp"))
	} else if jw.Claim("http://example.com/is_root") != true {
		t.Error(jw.Claim("http://example.com/is_root"))
	}
}

func TestJwt(t *testing.T) {
	jw := NewJwt()
	jw.SetHeader("alg", "none")
	jw.SetHeader("test", "test")
	jw.SetClaim("iss", "joe")
	jw.SetClaim("test", "test")

	if jw.Header("alg") != "none" {
		t.Error(jw.Header("alg"))
	} else if jw.Header("test") != "test" {
		t.Error(jw.Header("test"))
	} else if heads := jw.HeaderNames(); !heads["alg"] || !heads["test"] {
		t.Error(heads)
	} else if jw.Claim("iss") != "joe" {
		t.Error(jw.Claim("iss"))
	} else if jw.Claim("test") != "test" {
		t.Error(jw.Claim("test"))
	} else if clms := jw.ClaimNames(); !clms["iss"] || !clms["test"] {
		t.Error(clms)
	}

	jw.SetHeader("test", "")
	jw.SetClaim("test", nil)

	if jw.Header("alg") != "none" {
		t.Error(jw.Header("alg"))
	} else if jw.Header("test") != nil {
		t.Error(jw.Header("test"))
	} else if heads := jw.HeaderNames(); !heads["alg"] || heads["test"] {
		t.Error(heads)
	} else if jw.Claim("iss") != "joe" {
		t.Error(jw.Claim("iss"))
	} else if jw.Claim("test") != nil {
		t.Error(jw.Claim("test"))
	} else if clms := jw.ClaimNames(); !clms["iss"] || clms["test"] {
		t.Error(clms)
	}
}

func TestJwtEncode(t *testing.T) {
	jw := NewJwt()
	jw.SetHeader("alg", "none")
	jw.SetClaim("iss", "joe")
	jw.SetClaim("exp", 1300819380)
	jw.SetClaim("http://example.com/is_root", true)

	buff, err := jw.Encode()
	if err != nil {
		t.Fatal(err)
	}

	jw, err = ParseJwt(string(buff))
	if err != nil {
		t.Fatal(err, string(buff))
	}

	if jw.Header("alg") != "none" {
		t.Error(jw.Header("alg"))
	} else if jw.Claim("iss") != "joe" {
		t.Error(jw.Claim("iss"))
	} else if jw.Claim("exp") != 1300819380.0 {
		t.Error(jw.Claim("exp"))
	} else if jw.Claim("http://example.com/is_root") != true {
		t.Error(jw.Claim("http://example.com/is_root"))
	}
}

func TestToJson(t *testing.T) {
	jw := newJwt()
	jw.SetHeader("alg", "none")
	jw.SetClaim("iss", "joe")
	jw.SetClaim("exp", float64(1300819380))
	jw.SetClaim("http://example.com/is_root", true)

	headJs, clmsJs, err := jw.ToJson()
	if err != nil {
		t.Fatal(err)
	}

	var head map[string]interface{}
	if err := json.Unmarshal(headJs, &head); err != nil {
		t.Fatal(err, string(headJs))
	} else if !reflect.DeepEqual(head, jw.head) {
		t.Fatal(err, string(headJs))
	}

	var clms map[string]interface{}
	if err := json.Unmarshal(clmsJs, &clms); err != nil {
		t.Fatal(err, string(clmsJs))
	} else if !reflect.DeepEqual(clms, jw.clms) {
		t.Fatal(err, string(clmsJs))
	}
}
