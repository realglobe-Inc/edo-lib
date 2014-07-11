package util

import (
	"strings"
)

// path が apiPath 以下のパスかどうか調べる。
// path と apiPath の先頭は / であり、末尾は / ではないとする。
func HasApiPath(path, apiPath string) bool {
	if apiPath == "/" {
		return len(path) >= 2
	} else {
		return strings.HasPrefix(path, apiPath+"/")
	}
}

// hasPrefix が true だった path から apiPath 以降の部分を取り出す。
// ただし、返り値の先頭は常に / になるようにする。
func TrimApiPath(path, apiPath string) string {
	if apiPath == "/" {
		return path
	} else {
		return strings.TrimPrefix(path, apiPath)
	}
}
