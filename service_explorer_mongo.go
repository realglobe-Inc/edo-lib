package driver

import (
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// map[string]string から serviceExplorerTree をつくる。
func containerToServiceExplorerTree(data interface{}) (interface{}, error) {
	tree := newServiceExplorerTree()
	tree.fromContainer(data.(map[string]string))
	return tree, nil
}

// スレッドセーフ。
func NewMongoServiceExplorer(url, dbName, collName string, expiDur time.Duration) (ServiceExplorer, error) {
	base, err := newMongoKeyValueStore(url, dbName, collName, expiDur)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	base.MongoUnmarshal = containerToServiceExplorerTree
	base.MongoTake = containerMongoTake
	return newServiceExplorer(base), nil
}
