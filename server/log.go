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
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog"
	"github.com/realglobe-Inc/go-lib/rglog/level"
	"net/http"
	"net/http/httputil"
)

const logRoot = "github.com/realglobe-Inc"

var log = rglog.Logger(logRoot + "/edo-lib/server")

func LogRequest(lv level.Level, r *http.Request, useBody bool, prefs ...interface{}) {
	if log.IsLoggable(lv) {
		buff, err := httputil.DumpRequest(r, useBody)
		if err != nil {
			log.Warn(append(prefs, erro.Unwrap(err)))
			log.Debug(append(prefs, erro.Wrap(err)))
			return
		}
		log.Log(lv, append(prefs, "Request: ", string(buff))...)
	}
}

func LogResponse(lv level.Level, r *http.Response, useBody bool, prefs ...interface{}) {
	if log.IsLoggable(lv) {
		buff, err := httputil.DumpResponse(r, useBody)
		if err != nil {
			log.Warn(append(prefs, erro.Unwrap(err)))
			log.Debug(append(prefs, erro.Wrap(err)))
			return
		}
		log.Log(lv, append(prefs, "Response: ", string(buff))...)
	}
}
