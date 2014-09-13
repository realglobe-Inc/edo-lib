package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"path/filepath"
)

// 非キャッシュ用。
func NewFileLoginRegistry(path string) LoginRegistry {
	return newFileDriver(path)
}

func (reg *fileDriver) User(accToken string) (usrUuid string, err error) {
	path := filepath.Join(reg.path, accToken+".json")

	if err := readFromJson(path, &usrUuid); err != nil {
		return "", erro.Wrap(err)
	}
	return usrUuid, nil
}
