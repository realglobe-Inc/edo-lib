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

package server

import (
	"io"
	"net/http"

	"github.com/realglobe-Inc/go-lib/erro"
)

func CopyResponse(w http.ResponseWriter, resp *http.Response) error {
	// ヘッダフィールドのコピー。
	for key, vals := range resp.Header {
		w.Header().Set(key, vals[0])
		for _, val := range vals[1:] {
			w.Header().Add(key, val)
		}
	}

	// ステータスのコピー。
	w.WriteHeader(resp.StatusCode)

	// ボディのコピー。
	if _, err := io.Copy(w, resp.Body); err != nil {
		return erro.Wrap(err)
	}

	return nil
}
