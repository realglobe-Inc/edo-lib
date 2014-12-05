package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"net/http"
	"net/url"
	"regexp"
)

// URL を分解する。
// <scheme>://<host><remain>
func SplitUrl(rawUrl string) (scheme, host, remain string, err error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", "", "", erro.Wrap(err)
	}

	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}
	return u.Scheme, u.Host, path, nil
}

var slashes *regexp.Regexp

func init() {
	slashes = regexp.MustCompile("/+")
}

func MergeSlash(str string) string {
	return slashes.ReplaceAllString(str, "/")
}

func UrlPrefix(r *http.Request) string {
	var prefix string
	if s := r.Header.Get("X-Forwarded-Proto"); s != "" {
		prefix = s
	} else if s := r.Header.Get("X-Forwarded-Ssl"); s == "on" {
		prefix = "https"
	} else {
		prefix = "http"
	}

	prefix += "://"

	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		prefix += h
	} else {
		prefix += r.Host
	}

	return prefix
}
