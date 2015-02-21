package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"testing"
)

func TestHs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	type param struct {
		key interface{}
		crypto.Hash
	}
	for _, p := range []param{{test256Key, crypto.SHA256}, {test384Key, crypto.SHA384}, {test512Key, crypto.SHA512}} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := hsSign(p.key, p.Hash, plain); err != nil {
				t.Fatal(err)
			} else if err := hsVerify(p.key, p.Hash, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestRs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	for _, hGen := range []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := rsSign(testRsaKey, hGen, plain); err != nil {
				t.Fatal(err)
			} else if err := rsVerify(&testRsaKey.PublicKey, hGen, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestEs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	type param struct {
		key *ecdsa.PrivateKey
		crypto.Hash
	}
	for _, p := range []param{{testEcdsa256Key, crypto.SHA256}, {testEcdsa384Key, crypto.SHA384}, {testEcdsa521Key, crypto.SHA512}} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := esSign(p.key, p.Hash, plain); err != nil {
				t.Fatal(err)
			} else if err := esVerify(&p.key.PublicKey, p.Hash, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}

func TestPs(t *testing.T) {
	buff := []byte{}
	for ; len(buff) < 50; buff = append(buff, byte(len(buff))) {
	}
	for _, hGen := range []crypto.Hash{crypto.SHA256, crypto.SHA384, crypto.SHA512} {
		for i := 0; i < 50; i += 10 {
			plain := make([]byte, i)
			copy(plain, buff)
			if sig, err := psSign(testRsaKey, hGen, plain); err != nil {
				t.Fatal(err)
			} else if err := psVerify(&testRsaKey.PublicKey, hGen, sig, plain); err != nil {
				t.Error(sig)
				t.Fatal(err)
			}
		}
	}
}
