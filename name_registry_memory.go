package driver

import ()

// 非キャッシュ用。
type MemoryNameRegistry struct {
	tree *nameTree
}

func NewMemoryNameRegistry() *MemoryNameRegistry {
	return &MemoryNameRegistry{newNameTree()}
}

func (reg *MemoryNameRegistry) Address(name string) (addr string, err error) {
	return reg.tree.address(name), nil
}
func (reg *MemoryNameRegistry) Addresses(name string) (addrs []string, err error) {
	return reg.tree.addresses(name), nil
}
func (reg *MemoryNameRegistry) AddAddress(name, addr string) {
	reg.tree.add(name, addr)
}
func (reg *MemoryNameRegistry) RemoveAddress(name string) {
	reg.tree.remove(name)
}
