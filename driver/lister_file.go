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
	"github.com/realglobe-Inc/go-lib/erro"
	"os"
	"path/filepath"
	"time"
)

// キャッシュ用。
type fileLister struct {
	path      string
	pathToKey func(string) string
	staleDur  time.Duration
	expiDur   time.Duration
}

// スレッドセーフではない。
func newFileLister(path string, pathToKey func(string) string, staleDur, expiDur time.Duration) *fileLister {
	if pathToKey == nil {
		pathToKey = func(path string) string { return path }
	}
	return &fileLister{path, pathToKey, staleDur, expiDur}
}

func (drv *fileLister) getStamp(fi os.FileInfo) *Stamp {
	now := time.Now()
	stmp := getFileStamp(fi)
	stmp.StaleDate = now.Add(drv.staleDur)
	stmp.ExpiDate = now.Add(drv.expiDur)
	return stmp
}

func (drv *fileLister) Keys(caStmp *Stamp) (keys map[string]bool, newCaStmp *Stamp, err error) {
	fi, err := os.Stat(drv.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, erro.Wrap(err)
		}
	}

	newCaStmp = drv.getStamp(fi)

	// 対象のスタンプを取得。

	if caStmp != nil && !caStmp.Older(newCaStmp) {
		// 要求元のキャッシュより新しそうではなかった。
		return nil, newCaStmp, nil
	}

	// 要求元のキャッシュより新しそう。

	keys = map[string]bool{}
	if err := filepath.Walk(drv.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return erro.Wrap(err)
		} else if info.IsDir() {
			return nil
		}

		key := drv.pathToKey(path[len(drv.path)+1:]) // +1 は / の分。
		if key == "" {
			return nil
		}

		keys[key] = true
		return nil
	}); err != nil {
		return nil, nil, erro.Wrap(err)
	}

	return keys, newCaStmp, nil
}

func (drv *fileLister) Close() error {
	return nil
}