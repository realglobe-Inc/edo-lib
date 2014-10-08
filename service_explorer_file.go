package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

func serviceExplorerTreeMarshal(value interface{}) (data []byte, err error) {
	data, err = json.Marshal(value.(*serviceExplorerTree).toContainer())
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return data, nil
}

// data を JSON として、map[string]string にデコードしてから serviceExplorerTree をつくる。
func serviceExplorerTreeUnmarshal(data []byte) (interface{}, error) {
	var uriToUuid map[string]string
	if err := json.Unmarshal(data, &uriToUuid); err != nil {
		return nil, erro.Wrap(err)
	}

	tree := newServiceExplorerTree()
	tree.fromContainer(uriToUuid)
	return tree, nil
}

// スレッドセーフ。
func NewFileServiceExplorer(path string, expiDur time.Duration) ServiceExplorer {
	return newServiceExplorer(NewFileKeyValueStore(path, jsonKeyGen, serviceExplorerTreeMarshal, serviceExplorerTreeUnmarshal, expiDur))
}
