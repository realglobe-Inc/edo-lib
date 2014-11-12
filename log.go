package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

var log = rglog.Logger("github.com/realglobe-Inc/edo/util")

func initLog(root string, lv level.Level, key string, hndl handler.Handler) handler.Handler {
	rootLog := rglog.Logger(root)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)
	hndl.SetLevel(lv)
	rootLog.AddHandler(key, hndl)
	return hndl
}

func InitConsoleLog(root string) handler.Handler {
	return initLog(root, level.INFO, "console", handler.NewConsoleHandlerUsing(handler.LevelOnlyFormatter))
}

func initFileLog(root string, lv level.Level, path string) {
	initLog(root, lv, "file", handler.NewRotateHandler(path, 10*(1<<20), 10))
	log.Debug("Logging into file " + path + ".")
	return
}

func initFluentdLog(root string, lv level.Level, addr, tag string) {
	initLog(root, lv, "fluentd", handler.NewFluentdHandler(addr, tag))
	log.Debug("Logging into fluentd " + addr + ".")
	return
}

func SetupLog(root, logType string, logLv level.Level, logPath, fluAddr, fluTag string) error {
	switch logType {
	case "":
	case "file":
		initFileLog(root, logLv, logPath)
	case "fluentd":
		initFluentdLog(root, logLv, fluAddr, fluTag)
	default:
		return erro.New("invalid log type " + logType + ".")
	}
	return nil
}
