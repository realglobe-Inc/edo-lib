package log

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

var log = rglog.Logger("github.com/realglobe-Inc/edo/util/log")

const (
	TypeConsole = "console"
	TypeFile    = "file"
	TypeFluentd = "fluentd"
)

func setup(root string, lv level.Level, key string, hndl handler.Handler) handler.Handler {
	rootLog := rglog.Logger(root)
	if curLv := rootLog.Level(); curLv.Higher(lv) {
		rootLog.SetLevel(lv)
	}
	rootLog.SetUseParent(false)
	hndl.SetLevel(lv)
	if old := rootLog.AddHandler(key, hndl); old != nil {
		old.Close()
	}
	return hndl
}

func InitConsole(root string) {
	SetupConsole(root, level.INFO)
}

func SetupConsole(root string, lv level.Level) {
	if level.OFF.Higher(lv) {
		setup(root, lv, TypeConsole, handler.NewConsoleHandlerUsing(handler.LevelOnlyFormatter))
		log.Debug("Logging into console")
	}
	return
}

func SetupFile(root string, lv level.Level, path string, limit int64, num int) {
	if level.OFF.Higher(lv) {
		setup(root, lv, TypeFile, handler.NewRotateHandler(path, limit, num))
		log.Debug("Logging into file " + path)
	}
	return
}

func SetupFluentd(root string, lv level.Level, addr, tag string) {
	if level.OFF.Higher(lv) {
		setup(root, lv, TypeFluentd, handler.NewFluentdHandler(addr, tag))
		log.Debug("Logging into fluentd " + addr)
		return
	}
}

type FileOption interface {
	LogFilePath() string
	LogFileLimit() int64
	LogFileNumber() int
}

type FluentdOption interface {
	LogFluentdAddress() string
	LogFluentdTag() string
}

func Setup(root, logType string, logLv level.Level, opt interface{}) error {
	switch logType {
	case "":
	case TypeConsole:
		SetupConsole(root, logLv)
	case TypeFile:
		o, ok := opt.(FileOption)
		if !ok {
			return erro.New("log type " + logType + " requires option")
		}
		SetupFile(root, logLv, o.LogFilePath(), o.LogFileLimit(), o.LogFileNumber())
	case TypeFluentd:
		o, ok := opt.(FluentdOption)
		if !ok {
			return erro.New("log type " + logType + " requires option")
		}
		SetupFluentd(root, logLv, o.LogFluentdAddress(), o.LogFluentdTag())
	default:
		return erro.New("log type " + logType + " is unsupported")
	}
	return nil
}
