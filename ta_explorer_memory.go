package driver

import (
	"time"
)

type MemoryTaExplorer struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryTaExplorer(expiDur time.Duration) *MemoryTaExplorer {
	return &MemoryTaExplorer{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryTaExplorer) ServiceUuid(servUri string, caStmp *Stamp) (servUuid string, newCaStmp *Stamp, err error) {
	return (&taExplorer{reg.base}).ServiceUuid(servUri, caStmp)
}

func (reg *MemoryTaExplorer) SetServiceUuids(uriToUuid map[string]string) {
	tree := newTaExplorerTree()
	tree.fromContainer(uriToUuid)
	reg.base.Put("list", tree)
}
