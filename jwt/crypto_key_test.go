package jwt

import (
	"bytes"
	"testing"
)

func TestAes128Kw(t *testing.T) {
	// JWE Appendix A.3 より。
	key, err := KeyFromJwkMap(map[string]interface{}{
		"kty": "oct",
		"k":   "GawgguFyGrWKav7AX4VKUg",
	})
	if err != nil {
		t.Fatal(err)
	}
	plain := []byte{
		4, 211, 31, 197, 84, 157, 252, 254, 11, 100, 157, 250, 63, 170, 106, 206,
		107, 124, 212, 45, 111, 107, 9, 219, 200, 177, 0, 240, 143, 156, 44, 207,
	}
	encrypted := []byte{
		232, 160, 123, 211, 183, 76, 245, 132, 200, 128, 123, 75, 190, 216, 22, 67,
		201, 138, 193, 186, 9, 91, 122, 31, 246, 90, 28, 139, 57, 3, 76, 124,
		193, 11, 98, 37, 173, 61, 104, 57,
	}

	if e, err := encryptAesKw(key.([]byte), plain); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(e, encrypted) {
		t.Error(e)
		t.Error(encrypted)
	} else if p, err := decryptAesKw(key.([]byte), encrypted); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(p, plain) {
		t.Error(p)
		t.Error(plain)
	}
}

func TestAesKw(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 8*8; buff = append(buff, byte(len(buff))) {
	}

	for _, keyLen := range []int{16, 24, 32} {
		key := buff[:keyLen]

		for i := 2; i <= 8; i++ {
			plain := buff[:8*i]

			if encrypted, err := encryptAesKw(key, plain); err != nil {
				t.Fatal(err)
			} else if p, err := decryptAesKw(key, encrypted); err != nil {
				t.Fatal(err)
			} else if !bytes.Equal(p, plain) {
				t.Error(p)
				t.Error(plain)
			}
		}
	}
}
