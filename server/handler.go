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
	"html/template"
	"net/http"

	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

// デバッグログにリクエストボディを記録するかどうか。
var Debug = false

// WrapPage や WrapApi に渡す処理。
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// 処理がパニックやエラーで終わったら、適当なレスポンスを HTML で返す。
func WrapPage(stopper *Stopper, f HandlerFunc, errTmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logPref string

		// panic 対策。
		defer func() {
			if rcv := recover(); rcv != nil {
				RespondErrorHtml(w, r, erro.New(rcv), errTmpl, logPref)
				return
			}
		}()

		if stopper != nil {
			stopper.Stop()
			defer stopper.Unstop()
		}

		logPref = ParseSender(r) + ": "

		LogRequest(level.DEBUG, r, Debug, logPref)

		if err := f(w, r); err != nil {
			RespondErrorHtml(w, r, erro.Wrap(err), errTmpl, logPref)
			return
		}
	}
}

// 処理がパニックやエラーで終わったら、適当なレスポンスを JSON で返す。
func WrapApi(stopper *Stopper, f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logPref string

		// panic 対策。
		defer func() {
			if rcv := recover(); rcv != nil {
				RespondErrorJson(w, r, erro.New(rcv), logPref)
				return
			}
		}()

		if stopper != nil {
			stopper.Stop()
			defer stopper.Unstop()
		}

		logPref = ParseSender(r) + ": "

		LogRequest(level.DEBUG, r, Debug, logPref)

		if err := f(w, r); err != nil {
			RespondErrorJson(w, r, erro.Wrap(err), logPref)
			return
		}
	}
}
