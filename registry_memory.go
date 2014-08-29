package driver

import (
	"time"
)

// メモリ上のモックアップ。

// ログイン。
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

// JavaScript.
type MemoryJsRegistry struct {
	objs map[string]map[string]*Object
}

func NewMemoryJsRegistry() *MemoryJsRegistry {
	return &MemoryJsRegistry{map[string]map[string]*Object{}}
}

func (reg *MemoryJsRegistry) Object(dir, objName string) (*Object, error) {
	nameToObj := reg.objs[dir]
	if nameToObj == nil {
		return nil, nil
	}
	return nameToObj[objName], nil
}
func (reg *MemoryJsRegistry) AddObject(dir, objName string, obj *Object) error {
	nameToObj := reg.objs[dir]
	if nameToObj == nil {
		nameToObj = map[string]*Object{}
		reg.objs[dir] = nameToObj
	}
	nameToObj[objName] = obj
	return nil
}
func (reg *MemoryJsRegistry) RemoveObject(dir, objName string) error {
	nameToObj := reg.objs[dir]
	if nameToObj == nil {
		return nil
	}
	delete(nameToObj, objName)
	return nil
}

// ユーザー情報。
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

// ジョブ。
type MemoryJobRegistry struct {
	ress map[string]*JobResult
}

func NewMemoryJobRegistry() *MemoryJobRegistry {
	return &MemoryJobRegistry{map[string]*JobResult{}}
}

func (reg *MemoryJobRegistry) Result(jobId string) (res *JobResult, err error) {
	return reg.ress[jobId], nil
}
func (reg *MemoryJobRegistry) AddResult(jobId string, res *JobResult, deadline time.Time) error {
	reg.ress[jobId] = res
	return nil
}

// 別名。
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

// イベント。
type MemoryEventRegistry struct {
	trees map[string]*eventTree
}

func NewMemoryEventRegistry() *MemoryEventRegistry {
	return &MemoryEventRegistry{map[string]*eventTree{}}
}

func (reg *MemoryEventRegistry) Handler(usrUuid, event string) (Handler, error) {
	tree := reg.trees[usrUuid]
	if tree == nil {
		return nil, nil
	}
	return tree.handler(event), nil
}
func (reg *MemoryEventRegistry) AddHandler(usrUuid, event string, hndl Handler) error {
	tree := reg.trees[usrUuid]
	if tree == nil {
		tree = newEventTree()
		reg.trees[usrUuid] = tree
	}
	tree.add(event, hndl)
	return nil
}
func (reg *MemoryEventRegistry) RemoveHandler(usrUuid, event string) error {
	tree := reg.trees[usrUuid]
	if tree == nil {
		return nil
	}
	tree.remove(event)
	return nil
}

// サービス。
type MemoryServiceRegistry struct {
	*serviceTree
}

func NewMemoryServiceRegistry() *MemoryServiceRegistry {
	return &MemoryServiceRegistry{newServiceTree()}
}

func (reg *MemoryServiceRegistry) Service(endPt string) (servUuid string, err error) {
	return reg.service(endPt), nil
}
func (reg *MemoryServiceRegistry) AddService(endPt string, servUuid string) {
	reg.add(endPt, servUuid)
}
func (reg *MemoryServiceRegistry) RemoveService(endPt string) {
	reg.remove(endPt)
}

// ID プロバイダ。
type MemoryIdProviderRegistry struct {
	idps map[string]*IdProvider
}

func NewMemoryIdProviderRegistry() *MemoryIdProviderRegistry {
	return &MemoryIdProviderRegistry{map[string]*IdProvider{}}
}

func (reg *MemoryIdProviderRegistry) IdProviders() ([]*IdProvider, error) {
	idps := []*IdProvider{}
	for _, idp := range reg.idps {
		idps = append(idps, idp)
	}
	return idps, nil
}
func (reg *MemoryIdProviderRegistry) AddIdProvider(idp *IdProvider) {
	reg.idps[idp.IdpUuid] = idp
}
func (reg *MemoryIdProviderRegistry) RemoveIdProvider(idpUuid string) {
	delete(reg.idps, idpUuid)
}
