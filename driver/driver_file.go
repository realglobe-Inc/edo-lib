// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
