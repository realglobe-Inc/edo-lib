package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// 非キャッシュ用。
func NewFileNameRegistry(path string) NameRegistry {
	return newFileDriver(path)
}

func (reg *fileDriver) Address(name string) (addr string, err error) {
	path := filepath.Join(reg.path, name+".json")

	if err := readFromJson(path, &addr); err != nil { // 改行とかに煩わされないので JSON 文字列で。
		return "", erro.Wrap(err)
	}
	return addr, nil
}
func (reg *fileDriver) Addresses(name string) (addrs []string, err error) {
	cont := map[string]string{}

	fis, err := ioutil.ReadDir(reg.path)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		} else if !strings.HasSuffix(fi.Name(), ".json") {
			continue
		}

		curName := strings.TrimSuffix(fi.Name(), ".json")

		if !strings.HasSuffix(curName, name) {
			// 部分木以外はスルー。
			continue
		}

		path := filepath.Join(reg.path, fi.Name())

		var addr string
		if err := readFromJson(path, &addr); err != nil {
			return nil, erro.Wrap(err)
		}

		cont[curName] = addr
	}

	tree := newNameTree()
	tree.fromContainer(cont)

	return tree.addresses(name), nil
}
