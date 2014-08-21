package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
)

var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/edo/util")
}

func initLog(root string, lv level.Level, hndlGenerate func() (handler.Handler, error)) (handler.Handler, error) {
	rootLog := rglog.GetLogger(root)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)
	hndl, err := hndlGenerate()
	if err != nil {
		return nil, erro.Wrap(err)
	}
	hndl.SetLevel(lv)
	rootLog.AddHandler(hndl)
	return hndl, nil
}

func InitLog(root string) handler.Handler {
	hndl, _ := initLog(root, level.INFO, func() (handler.Handler, error) {
		return handler.NewConsoleHandlerUsing(handler.LevelOnlyFormatter), nil
	})
	return hndl
}

func InitFileLog(root string, lv level.Level, path string) error {
	if _, err := initLog(root, lv, func() (handler.Handler, error) {
		return handler.NewRotateHandler(path, 10*(1<<20), 10)
	}); err != nil {
		return erro.Wrap(err)
	}
	log.Debug("Logging into file " + path + ".")
	return nil
}

func InitFluentdLog(root string, lv level.Level, addr, tag string) error {
	if _, err := initLog(root, lv, func() (handler.Handler, error) {
		return handler.NewFluentdHandler(addr, tag)
	}); err != nil {
		return erro.Wrap(err)
	}
	log.Debug("Logging into fluentd " + addr + ".")
	return nil
}
