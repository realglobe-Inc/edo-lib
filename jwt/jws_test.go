package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func TestJws(t *testing.T) {
	// JSON Web Token (JWT) より。
	m := map[string]interface{}{
		"kty": "RSA",
		"n":   "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
		"e":   "AQAB",
	}
	_, key, err := PublicKeyFromJwkMap(m)
	if err != nil {
		t.Fatal(err)
	}
	keySet := map[string]crypto.PublicKey{"": key}

	raw := "eyJhbGciOiJSUzI1NiJ9" +
		"." + "eyJpc3MiOiJqb2UiLA0KICJleHAiOjEzMDA4MTkzODAsDQogImh0dHA6Ly9leGFtcGxlLmNvbS9pc19yb290Ijp0cnVlfQ" +
		"." + "cC4hiUPoj9Eetdgtv3hF80EGrhuB__dzERat0XF9g2VtQgr9PJbu3XOiZj5RZmh7AAuHIm4Bh-0Qc_lF5YKt_O8W2Fp5jujGbds9uJdbF9CUAr7t1dnZcAcQjbKBYNX4BAynRFdiuB--f_nZLgrnbyTyWzO75vRK5h6xBArLIARNPvkSjtQBMHlb1L07Qe7K0GarZRmB_eSN9383LcOLn6_dO--xi12jzDwusC-eOkHWEsqtFZESc6BfI7noOPqvhJ1phCnvWh6IeYI2w9QOYEUipUTI8np6LbgGY9Fs98rqVt5AXLIhWkWywlVmtVrBp0igcN_IoypGlUPQGe77Rw"

	s, err := ParseJws(raw)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Verify(keySet); err != nil {
		t.Fatal(err)
	}
}

func TestJwsSignAndVerify(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 1152)
	if err != nil {
		t.Fatal(err)
	}
	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	priKeySet := map[string]crypto.PrivateKey{
		"none":  nil,
		"RS256": rsaKey,
		"RS384": rsaKey,
		"RS512": rsaKey,
		"ES256": ecdsaKey,
		"ES384": ecdsaKey,
		"ES512": ecdsaKey,
		"PS256": rsaKey,
		"PS384": rsaKey,
		"PS512": rsaKey,
	}
	pubKeySet := map[string]crypto.PublicKey{
		"none":  nil,
		"RS256": &rsaKey.PublicKey,
		"RS384": &rsaKey.PublicKey,
		"RS512": &rsaKey.PublicKey,
		"ES256": &ecdsaKey.PublicKey,
		"ES384": &ecdsaKey.PublicKey,
		"ES512": &ecdsaKey.PublicKey,
		"PS256": &rsaKey.PublicKey,
		"PS384": &rsaKey.PublicKey,
		"PS512": &rsaKey.PublicKey,
	}

	for _, alg := range []string{"none", "RS256", "RS384", "RS512", "ES256", "ES384", "ES512", "PS256", "PS384", "PS512"} {
		jw := NewJws()
		jw.SetHeader("alg", alg)
		jw.SetHeader("kid", alg)
		jw.SetClaim("iss", "joe")
		jw.SetClaim("exp", float64(1300819380))
		jw.SetClaim("http://example.com/is_root", true)

		if err := jw.Sign(priKeySet); err != nil {
			b, _ := jw.Encode()
			h, c, _ := jw.ToJson()
			t.Fatal(err, alg, string(b), string(h), string(c))
		} else if err := jw.Verify(pubKeySet); err != nil {
			b, _ := jw.Encode()
			h, c, _ := jw.ToJson()
			t.Fatal(err, alg, string(b), string(h), string(c))
		}
	}
}