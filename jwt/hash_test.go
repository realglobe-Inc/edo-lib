package jwt

import (
	"testing"
)

func TestHashFunction(t *testing.T) {
	for _, alg := range testAlgList {
		h, err := HashFunction(alg)
		if err != nil {
			t.Fatal(err)
		}
		if h == nil {
			if alg != "none" {
				t.Error(alg, h)
			}
		}
	}
}
