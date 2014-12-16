package driver

import (
	"os"
	"strconv"
)

const (
	dirPerm  = 0755
	filePerm = 0644
)

// ダイジェストはタイムスタンプの秒未満にファイルサイズを足して 16 進数で表した文字列。
func getFileStamp(fi os.FileInfo) *Stamp {
	date := fi.ModTime()
	return &Stamp{
		Date:   date,
		Digest: strconv.FormatInt(int64(date.Nanosecond())+fi.Size(), 16),
	}
}
