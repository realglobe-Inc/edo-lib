package util

import (
	"testing"
)

func TestSecureRandomBytes(t *testing.T) {
	for i := 0; i < 100; i++ {
		buff, err := SecureRandomBytes(i)
		if err != nil {
			t.Fatal(err)
		} else if len(buff) != i {
			t.Error(i, len(buff))
		}
	}
}
