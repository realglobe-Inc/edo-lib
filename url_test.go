package util

import (
	"testing"
)

func TestSplitUrl(t *testing.T) {
	for _, s := range []string{"http", "https"} {
		for _, h := range []string{"example.com", "example.com:8080"} {
			for _, r := range []string{"/", "/index.html", "/a/index.html", "/index.html?a=b"} {
				rawUrl := s + "://" + h + r
				scheme, host, remain, err := SplitUrl(rawUrl)
				if err != nil {
					t.Fatal(err)
				} else if scheme != s {
					t.Error(scheme+" "+host+" "+remain, rawUrl)
				} else if host != h {
					t.Error(scheme+" "+host+" "+remain, rawUrl)
				} else if remain != r {
					t.Error(scheme+" "+host+" "+remain, rawUrl)
				}
			}
		}
	}
}
