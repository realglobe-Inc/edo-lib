package driver

import (
	"time"
)

// nameTree を保存する。
type MemoryNameRegistry struct {
	base KeyValueStore
}

// スレッドセーフ。
func NewMemoryNameRegistry(expiDur time.Duration) *MemoryNameRegistry {
	return &MemoryNameRegistry{NewMemoryKeyValueStore(expiDur)}
}

func (reg *MemoryNameRegistry) Address(name string, caStmp *Stamp) (addr string, newCaStmp *Stamp, err error) {
	return (&nameRegistry{reg.base}).Address(name, caStmp)
}
func (reg *MemoryNameRegistry) Addresses(name string, caStmp *Stamp) (addrs []string, newCaStmp *Stamp, err error) {
	return (&nameRegistry{reg.base}).Addresses(name, caStmp)
}

func (reg *MemoryNameRegistry) SetAddresses(cont map[string]string) {
	tree := newNameTree()
	tree.fromContainer(cont)
	reg.base.Put("names", tree)
}
