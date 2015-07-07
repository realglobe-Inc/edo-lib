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

// ログ関係。
package log

import (
	"github.com/realglobe-Inc/go-lib/erro"
	"github.com/realglobe-Inc/go-lib/rglog"
	"github.com/realglobe-Inc/go-lib/rglog/handler"
	"github.com/realglobe-Inc/go-lib/rglog/level"
)

const logRoot = "github.com/realglobe-Inc"

var log = rglog.Logger(logRoot + "/edo-lib/log")

const (
	TypeConsole = "console"
	TypeFile    = "file"
	TypeFluentd = "fluentd"
)

func setup(root string, lv level.Level, key string, hndl handler.Handler) {
	rootLog := rglog.Logger(root)
	if curLv := rootLog.Level(); curLv.Higher(lv) {
		rootLog.SetLevel(lv)
	}
	rootLog.SetUseParent(false)

	if hndl == nil {
		if old := rootLog.RemoveHandler(key); old != nil {
			old.Close()
		}
		return
	}

	hndl.SetLevel(lv)
	if old := rootLog.AddHandler(key, hndl); old != nil {
		old.Close()
	}
	return
}

func InitConsole(root string) {
	SetupConsole(root, level.INFO)
}

func SetupConsole(root string, lv level.Level) {
	var hndl handler.Handler
	if level.OFF.Higher(lv) {
		hndl = handler.NewConsoleHandlerUsing(handler.LevelOnlyFormatter)
	}
	setup(root, lv, TypeConsole, hndl)
	log.Debug(lv, " logging into console")
	return
}

func SetupFile(root string, lv level.Level, path string, limit int64, num int) {
	var hndl handler.Handler
	if level.OFF.Higher(lv) {
		hndl = handler.NewRotateHandler(path, limit, num)
	}
	setup(root, lv, TypeFile, hndl)
	log.Debug(lv, " logging into file "+path)
	return
}

func SetupFluentd(root string, lv level.Level, addr, tag string) {
	var hndl handler.Handler
	if level.OFF.Higher(lv) {
		hndl = handler.NewFluentdHandler(addr, tag)
	}
	setup(root, lv, TypeFluentd, hndl)
	log.Debug(lv, " logging into fluentd "+addr)
	return
}

type FileParameter interface {
	LogFilePath() string
	LogFileLimit() int64
	LogFileNumber() int
}

type FluentdParameter interface {
	LogFluentdAddress() string
	LogFluentdTag() string
}

func Setup(root, logType string, logLv level.Level, param interface{}) error {
	switch logType {
	case "":
	case TypeConsole:
		SetupConsole(root, logLv)
	case TypeFile:
		p, ok := param.(FileParameter)
		if !ok {
			return erro.New("log type " + logType + " requires option")
		}
		SetupFile(root, logLv, p.LogFilePath(), p.LogFileLimit(), p.LogFileNumber())
	case TypeFluentd:
		p, ok := param.(FluentdParameter)
		if !ok {
			return erro.New("log type " + logType + " requires option")
		}
		SetupFluentd(root, logLv, p.LogFluentdAddress(), p.LogFluentdTag())
	default:
		return erro.New("log type " + logType + " is unsupported")
	}
	return nil
}
