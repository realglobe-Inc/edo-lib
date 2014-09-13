package driver

import ()

// 非キャッシュ用。
type MemoryUserRegistry struct {
	attrs map[string]map[string]interface{}
}

func NewMemoryUserRegistry() *MemoryUserRegistry {
	return &MemoryUserRegistry{map[string]map[string]interface{}{}}
}

func (reg *MemoryUserRegistry) Attributes(usrUuid string) (map[string]interface{}, error) {
	return reg.attrs[usrUuid], nil
}
func (reg *MemoryUserRegistry) AddAttributes(usrUuid string, attrs map[string]interface{}) {
	reg.attrs[usrUuid] = attrs
}
func (reg *MemoryUserRegistry) RemoveAttributes(usrUuid string) {
	delete(reg.attrs, usrUuid)
}
func (reg *MemoryUserRegistry) Attribute(usrUuid, attrName string) (interface{}, error) {
	attrs := reg.attrs[usrUuid]
	if attrs == nil {
		return nil, nil
	}
	return attrs[attrName], nil
}
func (reg *MemoryUserRegistry) AddAttribute(usrUuid, attrName string, attr interface{}) error {
	attrs := reg.attrs[usrUuid]
	if attrs == nil {
		attrs = map[string]interface{}{}
		reg.attrs[usrUuid] = attrs
	}
	attrs[attrName] = attr
	return nil
}
func (reg *MemoryUserRegistry) RemoveAttribute(usrUuid, attrName string) error {
	attrs := reg.attrs[usrUuid]
	if attrs == nil {
		return nil
	}
	delete(attrs, attrName)
	return nil
}
