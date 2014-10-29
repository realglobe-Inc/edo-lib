package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// map[string]string から taExplorerTree をつくる。
func containerToTaExplorerTree(data interface{}) (interface{}, error) {
	tree := newTaExplorerTree()
	tree.fromContainer(data.(map[string]string))
	return tree, nil
}

// スレッドセーフ。
func NewMongoTaExplorer(url, dbName, collName string, expiDur time.Duration) (TaExplorer, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.MongoUnmarshal = containerToTaExplorerTree
	base.MongoTake = containerMongoTake
	return newTaExplorer(base), nil
}
