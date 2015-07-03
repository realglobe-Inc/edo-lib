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
	"encoding/json"
	"github.com/realglobe-Inc/go-lib/erro"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestRespondErrorJson(t *testing.T) {
	origErr := NewError(http.StatusBadRequest, "invalid request", nil)

	w := httptest.NewRecorder()
	RespondErrorJson(w, nil, erro.Wrap(origErr))

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "application/json" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("application/json")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	data, _ := ioutil.ReadAll(w.Body)
	var buff struct {
		Status  int
		Message string
	}
	if err := json.Unmarshal(data, &buff); err != nil {
		t.Fatal(err)
	} else if buff.Status != origErr.Status() {
		t.Error(buff.Status)
		t.Fatal(origErr.Status())
	} else if buff.Message != origErr.Message() {
		t.Error(buff.Message)
		t.Fatal(origErr.Message())
	}
}

func TestRespondErrorHtml(t *testing.T) {
	origErr := NewError(http.StatusBadRequest, "invalid request", nil)

	w := httptest.NewRecorder()
	RespondErrorHtml(w, nil, erro.Wrap(origErr), nil)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}
}

func TestRespondErrorHtmlTemplate(t *testing.T) {
	origErr := NewError(http.StatusBadRequest, "invalid request", nil)

	file, err := ioutil.TempFile("", "edo-lib")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write([]byte("{{.Status}}")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	tmpl, err := template.ParseFiles(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RespondErrorHtml(w, nil, erro.Wrap(origErr), tmpl)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	buff, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	} else if string(buff) != strconv.Itoa(origErr.Status()) {
		t.Error(string(buff))
		t.Fatal(origErr.Status())
	}
}

func TestRespondErrorHtmlTemplateFunction(t *testing.T) {
	origErr := NewError(http.StatusBadRequest, "invalid request", erro.New("mazui"))

	file, err := ioutil.TempFile("", "edo-lib")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write([]byte("{{.Status}} {{.StatusText}} {{.Error}} {{.Debug}}")); err != nil {
		t.Fatal(err)
	}
	file.Close()

	tmpl, err := template.ParseFiles(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	RespondErrorHtml(w, nil, erro.Wrap(origErr), tmpl)

	if w.Code != origErr.Status() {
		t.Error(w.Code)
		t.Fatal(origErr.Status())
	} else if w.HeaderMap.Get("Content-Type") != "text/html" {
		t.Error(w.HeaderMap.Get("Content-Type"))
		t.Fatal("text/html")
	} else if w.Body == nil {
		t.Fatal("no body")
	}

	buff, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	} else if parts := strings.Fields(string(buff)); parts[0] != strconv.Itoa(origErr.Status()) {
		t.Error(string(buff))
		t.Fatal(origErr.Status())
	}
}
