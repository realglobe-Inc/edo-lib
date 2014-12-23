package util

import (
	"testing"
)

func TestSecureRandomString(t *testing.T) {
	for i := 0; i < 100; i++ {
		buff, err := SecureRandomString(i)
		if err != nil {
			t.Fatal(err)
		} else if len(buff) != i {
			t.Error(i, len(buff), " "+buff)
		} else if len(buff) > 0 && buff[len(buff)-1] == '=' {
			t.Error(i, len(buff), " "+buff)
		}
	}
}
