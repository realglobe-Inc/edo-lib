package driver

import ()

// 非キャッシュ用。
type MemoryLoginRegistry struct {
	usrs map[string]string
}

func NewMemoryLoginRegistry() *MemoryLoginRegistry {
	return &MemoryLoginRegistry{map[string]string{}}
}

func (reg *MemoryLoginRegistry) User(accToken string) (addr string, err error) {
	return reg.usrs[accToken], nil
}
func (reg *MemoryLoginRegistry) AddUser(accToken string, addr string) {
	reg.usrs[accToken] = addr
}
func (reg *MemoryLoginRegistry) RemoveUser(accToken string) {
	delete(reg.usrs, accToken)
}
