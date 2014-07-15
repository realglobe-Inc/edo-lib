package util

import (
	"testing"
)

func TestSplitUrl(t *testing.T) {
	scheme, host, remain, err := SplitUrl("https://localhost:8000/index.html")
	if err != nil {
		t.Fatal(err)
	} else if scheme != "https" {
		t.Error(scheme + " " + host + " " + remain)
	} else if host != "localhost:8000" {
		t.Error(scheme + " " + host + " " + remain)
	} else if remain != "/index.html" {
		t.Error(scheme + " " + host + " " + remain)
	}

	scheme, host, remain, err = SplitUrl("https://@c.b.d")
	if err != nil {
		t.Fatal(err)
	} else if scheme != "https" {
		t.Error(scheme + " " + host + " " + remain)
	} else if host != "@c.b.d" {
		t.Error(scheme + " " + host + " " + remain)
	} else if remain != "" {
		t.Error(scheme + " " + host + " " + remain)
	}
}
