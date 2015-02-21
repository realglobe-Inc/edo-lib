package jwt

import (
	"bytes"
	"strings"
	"testing"
)

func TestBase64Url(t *testing.T) {
	for b := []byte{}; len(b) < 100; b = append(b, byte(len(b))) {
		if s := base64UrlEncodeToString(b); strings.Index(s, "=") > 0 {
			t.Error(b)
			t.Error(s)
		} else if b2, err := base64UrlDecodeString(s); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(b2, b) {
			t.Error(b2)
			t.Error(b)
			t.Error(s)
		}
	}
}
