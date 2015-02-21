package jwt

import (
	"bytes"
	"testing"
)

func TestDef(t *testing.T) {
	for plain := []byte{}; len(plain) < 100; plain = append(plain, byte(len(plain))) {
		if d, err := defCompress(plain); err != nil {
			t.Fatal(err)
		} else if b2, err := defDecompress(d); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(b2, plain) {
			t.Error(b2)
			t.Error(plain)
		}
	}
}
