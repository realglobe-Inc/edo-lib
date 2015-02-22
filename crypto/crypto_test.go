package crypto

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsePemRsaPublic(t *testing.T) {
	key, err := ParsePem([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAr8+0/TRdzgwkHyB8DOWd
IiRSpwfs6JdlPrjDwAOpXQquwN36UDFHKtyQBeV8dw0t1imKwvtUFAmQbgQtcJ+3
GF3PDoDg/v5UQgxuUI/vNKiYG1BDjuwcXDnUT9fWCIqXy34M9z/HJFs3+BlmspZw
fOTULuH6wuqh64d+Pctn4srm4ZnJ93TQh8LtMXoBUiOJYoRXdl8sRmd5bO0BEP86
06ZjNC0F97I94sqYtrqhTvMgQlfmKkbsGmvJq6bbfPHtgEMH2KBDqTdXWaipdCoU
qt5O2Y0HU53Xh1T/I9hJL3EanSOtvY81qijVkUBVmKlfWO3X1+MXn0F3Hev3ZxMm
VwIDAQAB
-----END PUBLIC KEY-----`))
	if err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*rsa.PublicKey); !ok {
		t.Error(key)
	}
}

func TestParsePemRsaPrivate(t *testing.T) {
	key, err := ParsePem([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAr8+0/TRdzgwkHyB8DOWdIiRSpwfs6JdlPrjDwAOpXQquwN36
UDFHKtyQBeV8dw0t1imKwvtUFAmQbgQtcJ+3GF3PDoDg/v5UQgxuUI/vNKiYG1BD
juwcXDnUT9fWCIqXy34M9z/HJFs3+BlmspZwfOTULuH6wuqh64d+Pctn4srm4ZnJ
93TQh8LtMXoBUiOJYoRXdl8sRmd5bO0BEP8606ZjNC0F97I94sqYtrqhTvMgQlfm
KkbsGmvJq6bbfPHtgEMH2KBDqTdXWaipdCoUqt5O2Y0HU53Xh1T/I9hJL3EanSOt
vY81qijVkUBVmKlfWO3X1+MXn0F3Hev3ZxMmVwIDAQABAoIBAD3kEePdLnSdw42N
ov3rSyC8xrf0S0sqGMM2yfprj5CodRKCUl8uqc4F7VGWEvXaFtvArg+r6FJRd52z
LMfsAcm7JGwHpK0/nSvPMnp74QqZm0pqPA4xQl6ZIQumgLtrBrrlOe1Eb3d2AUL+
ti+CVEEzURrcBKnfbXb7sM5SL9UfYTvdznNlkZ4ZVWdjMX/xnzziTl7fbX5ijHGr
Vko7Z+iHwKqeKd2G2Fj0vRciC7l9evvHhaxQl3kW8TXM9AIljof84w5nrB7shZDX
6kaYnzI8mTCMDMyWLG/OWDkpx0heY6WpMyX18K+sPwod23zs1rHMH6K9Y82knUmE
962fWQkCgYEA1xXAlpwOODylqj//vsUnFAj7FlkevdUvOe/MQvWpDQVRNsDKLPtb
v2ThkQ0rGZp37tcpxneZjqGyghPf2WOQmQbel1iuOhEOImFoSImFVvaxAQRYx6m4
3zYWu1a9jN/F+nYmhXlFQqIw8KwN9TEMTojODZ/oRfMt3pRlAEhXvT0CgYEA0UFm
o968JXOdrAnODGkzGuYvcGSiunmyXQY+745fxouuhJVl1pTNeAC2Ct7dflKiAJIF
tN2gqY/G6d3FMMeYVlrASwwD/Dpbma+LQRkgZjCW8jAXGzSkMLKA8QSeXVk0JvQl
3DY08rxe7ORBN3CnPr5+jr6hsxc6TNOUjlI30yMCgYBemAN+eZ3LX+jgSotYxG6e
YiDDwGhDxvmhOnSUUmSKBHemY/3G8Ll2IJEP8UGuXgA3O8v0rG8NitHuYX1Gp4JV
uu60k1z0zsFvn3V0yX6qM46/SsEc9ukGykwPEmQFC/mPYN0qQJ6UYq6xeooc9vhZ
pdMxrM1DzmKzDIKrMCXeZQKBgQCJRvCA/LRNlYWQwXXtam5ebTgd8cdXslKy+E/9
dFectzIsRJ5koYYR/dVvWDnSj388BI+90c9+rZX/AsBEegyUSkDwetd6dwZ00lb2
w/cfUy0TgT0HWgeE8vXoJ/GEp+qwy2azCtS9kZpsqmmmZz8wyGPaXXFTPh+/GubQ
X1vEJQKBgDl7DijE1/jr2x9ACbgbtMrRk4HekzdrwOGgZ4TdsWgKfi/WWxbO6g5R
cGX8XbmB21BqExzN7jPhbvhtLi271jXiirdLInOYWx0/nn8Ks2qkdN5eE+XIqaeN
PC1uutoixe1WZTzrWYPIOFBXeQVFlUbnmZdj0LnqAJsIz1Vec9K8
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*rsa.PrivateKey); !ok {
		t.Error(key)
	}
}

func TestParsePemEcdsaPublic(t *testing.T) {
	key, err := ParsePem([]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEOCbnPn2SPA92u2G09XmrB9rTeqWv
SFeYEjDv3p7hDnDS+vrPmEQ3twGw7vn38JoIIhYdowJX4+deWcezFDtI1A==
-----END PUBLIC KEY-----`))
	if err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*ecdsa.PublicKey); !ok {
		t.Error(key)
	}
}

func TestParsePemEcdsaPrivate(t *testing.T) {
	key, err := ParsePem([]byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIEqu+CBiCePKkS6V1YLsjMsiEk86fV18cEHMgt0qLSwFoAoGCCqGSM49
AwEHoUQDQgAEOCbnPn2SPA92u2G09XmrB9rTeqWvSFeYEjDv3p7hDnDS+vrPmEQ3
twGw7vn38JoIIhYdowJX4+deWcezFDtI1A==
-----END EC PRIVATE KEY-----`))
	if err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*ecdsa.PrivateKey); !ok {
		t.Error(key)
	}
}

const testLabel = "edo-test"

func TestReadPem(t *testing.T) {
	f, err := ioutil.TempFile("", testLabel)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(f.Name())
	if _, err := f.Write([]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEOCbnPn2SPA92u2G09XmrB9rTeqWv
SFeYEjDv3p7hDnDS+vrPmEQ3twGw7vn38JoIIhYdowJX4+deWcezFDtI1A==
-----END PUBLIC KEY-----`)); err != nil {
		t.Fatal(err)
	}

	key, err := ReadPem(f.Name())
	if err != nil {
		t.Fatal(err)
	} else if _, ok := key.(*ecdsa.PublicKey); !ok {
		t.Error(key)
	}
}
