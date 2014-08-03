package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"regexp"
	"strings"
)

// URL を分解する。
// <scheme>://<host><remain>
func SplitUrl(url string) (scheme, host, remain string, err error) {
	idx := strings.Index(url, "://")
	if idx < 0 {
		return "", "", "", erro.New("invalid url " + url + ".")
	}

	scheme = url[:idx]
	host = url[idx+len("://"):]

	idx = strings.Index(host, "/")
	if idx >= 0 {
		remain = host[idx:]
		host = host[:idx]
	}

	return scheme, host, remain, nil
}

var slashes *regexp.Regexp

func init() {
	slashes = regexp.MustCompile("/+")
}

func MergeSlash(str string) string {
	return slashes.ReplaceAllString(str, "/")
}
