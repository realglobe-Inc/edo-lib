package util

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"github.com/realglobe-Inc/go-lib-rg/rglog"
	"github.com/realglobe-Inc/go-lib-rg/rglog/handler"
	"github.com/realglobe-Inc/go-lib-rg/rglog/level"
	"os"
	"path/filepath"
)

var log rglog.Logger

func init() {
	log = rglog.GetLogger("github.com/realglobe-Inc/edo/util")
}

const dirPerm = 0755

func InitLog(root string) handler.Handler {
	rootLog := rglog.GetLogger(root)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)
	hndl := handler.NewConsoleHandlerUsing(handler.LevelOnlyFormatter)
	hndl.SetLevel(level.INFO)
	rootLog.AddHandler(hndl)
	return hndl
}

func InitFileLog(root string, lv level.Level, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return erro.Wrap(err)
	}
	rootLog := rglog.GetLogger(root)
	rootLog.SetLevel(level.ALL)
	rootLog.SetUseParent(false)
	hndl, err := handler.NewRotateHandler(path, 10*(1<<20), 10)
	if err != nil {
		return erro.Wrap(err)
	}
	hndl.SetLevel(lv)
	rootLog.AddHandler(hndl)
	return nil
}
