package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

func nameTreeMarshal(value interface{}) (data []byte, err error) {
	data, err = json.Marshal(value.(*nameTree).toContainer())
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return data, nil
}

// data を JSON として、map[string]string にデータ型にデコードしてから nameTree をつくる。
func nameTreeUnmarshal(data []byte) (interface{}, error) {
	var cont map[string]string
	if err := json.Unmarshal(data, &cont); err != nil {
		return nil, erro.Wrap(err)
	}

	tree := newNameTree()
	tree.fromContainer(cont)
	return tree, nil
}

// スレッドセーフ。
func NewFileNameRegistry(path string, expiDur time.Duration) NameRegistry {
	return newNameRegistry(NewFileKeyValueStore(path, jsonKeyGen, nameTreeMarshal, nameTreeUnmarshal, expiDur))
}
