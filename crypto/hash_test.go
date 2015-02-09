package crypto

import (
	"crypto"
	"testing"
)

func TestHashFunction(t *testing.T) {
	for _, h := range []crypto.Hash{
		crypto.MD4,
		crypto.MD5,
		crypto.SHA1,
		crypto.SHA224,
		crypto.SHA256,
		crypto.SHA384,
		crypto.SHA512,
		crypto.MD5SHA1,
		crypto.RIPEMD160,
		crypto.SHA3_224,
		crypto.SHA3_256,
		crypto.SHA3_384,
		crypto.SHA3_512,
	} {
		s := HashFunctionString(h)
		h2, err := ParseHashFunction(s)
		if err != nil {
			t.Fatal(err)
		} else if h2 != h {
			t.Error(h2, h)
		}
	}
}

func TestParseUnknownHashFunction(t *testing.T) {
	_, err := ParseHashFunction("unknown")
	if err == nil {
		t.Error("no error")
	}
}

func TestUnknownHashFunctionString(t *testing.T) {
	str := HashFunctionString(crypto.Hash(1000000))
	if _, ok := strToHash[str]; ok {
		t.Error(str)
	}
}
