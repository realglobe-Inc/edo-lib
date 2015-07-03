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
	"bytes"
	"encoding/json"
	jsonutil "github.com/realglobe-Inc/edo-lib/json"
	"github.com/realglobe-Inc/go-lib/erro"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// JSON でエラーを返す。
func RespondErrorJson(w http.ResponseWriter, r *http.Request, origErr error, logPrefs ...interface{}) {
	e := ErrorFrom(origErr)
	log.Err(append(logPrefs, e.Status(), " "+http.StatusText(e.Status())+": "+e.Message())...)
	log.Debug(append(logPrefs, origErr)...)

	buff, err := json.Marshal(map[string]interface{}{
		tagStatus:  e.Status(),
		tagMessage: e.Message(),
	})
	if err != nil {
		log.Err(append(logPrefs, erro.Unwrap(err))...)
		log.Debug(append(logPrefs, erro.Wrap(err))...)
		// 最後の手段。たぶん正しい変換。
		buff = []byte(`{` +
			tagStatus + `=` + strconv.Itoa(e.Status()) + `,` +
			tagMessage + `="` + jsonutil.StringEscape(e.Message()) +
			`"}`)
	}

	w.Header().Set(tagContent_type, contTypeJson)
	w.WriteHeader(e.Status())
	if _, err := w.Write(buff); err != nil {
		log.Err(append(logPrefs, erro.Wrap(err))...)
	}
}

// HTML でエラーを返す。
// テンプレートでは以下が使える。
// {{.Status}}: HTTP ステータスコード。404 とか
// {{.StatusText}}: HTTP ステータスコード。Not Found とか
// {{.Error}}: エラー内容
// {{.Debug}}: エラー詳細
func RespondErrorHtml(w http.ResponseWriter, r *http.Request, origErr error, errTmpl *template.Template, logPrefs ...interface{}) {
	e := ErrorFrom(origErr)
	log.Err(append(logPrefs, e.Status(), " "+http.StatusText(e.Status())+": "+e.Message())...)
	log.Debug(append(logPrefs, origErr)...)

	var data []byte
	if errTmpl != nil {
		// テンプレートからユーザー向けの HTML をつくる。
		buff := &bytes.Buffer{}
		if err := errTmpl.Execute(buff, &templateData{base: e}); err != nil {
			log.Warn(append(logPrefs, erro.Unwrap(err))...)
			log.Debug(append(logPrefs, erro.Wrap(err))...)
		} else {
			data = buff.Bytes()
		}
	}

	if data == nil {
		// 自前でユーザー向けの HTML をつくる。
		msg := `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Error</title></head><body><h1>`
		msg += strconv.Itoa(e.Status()) + " " + http.StatusText(e.Status())
		msg += `</h1><p><b>` + tagMessage + `:</b> `
		msg += strings.Replace(html.EscapeString(e.Message()), "\n", "<br/>", -1)
		msg += `<br/><b>` + tagDebug + `:</b> `
		msg += strings.Replace(html.EscapeString(e.Error()), "\n", "<br/>", -1)
		msg += `</p></body></html>`
		data = []byte(msg)
	}

	w.Header().Set(tagContent_type, contTypeHtml)
	w.WriteHeader(e.Status())
	if _, err := w.Write(data); err != nil {
		log.Err(append(logPrefs, erro.Wrap(err))...)
	}
}

// テンプレートデータ。
type templateData struct {
	base    *Error
	statTxt string
	debug   template.HTML
}

func (this *templateData) Status() int {
	return this.base.Status()
}

func (this *templateData) StatusText() string {
	if this.statTxt == "" {
		this.statTxt = http.StatusText(this.base.Status())
	}
	return this.statTxt
}

func (this *templateData) Error() string {
	return this.base.Message()
}

func (this *templateData) Debug() template.HTML {
	if this.debug == "" {
		debug := template.HTMLEscapeString(this.base.Error())
		debug = strings.Replace(debug, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		debug = strings.Replace(debug, "\n", "<br/>", -1)
		this.debug = template.HTML(debug)
	}
	return this.debug
}
