package driver

import (
	"time"
)

type MemoryServiceExplorer struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryServiceExplorer(expiDur time.Duration) *MemoryServiceExplorer {
	return &MemoryServiceExplorer{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryServiceExplorer) ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	return (&serviceExplorer{reg.base}).ServiceUuid(servUri, caStmp)
}

func (reg *MemoryServiceExplorer) SetServiceUuids(uriToUuid map[string]string) {
	tree := newServiceExplorerTree()
	tree.fromContainer(uriToUuid)
	reg.base.Put("list", tree)
}
