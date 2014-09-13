package driver

import (
	"fmt"
	"strconv"
	"time"
)

// 非キャッシュ用。
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

// キャッシュ用。
type MemoryJsBackendRegistry struct {
	*MemoryJsRegistry
	stmps   map[string]map[string]*Stamp
	expiDur time.Duration
}

func NewMemoryJsBackendRegistry(expiDur time.Duration) *MemoryJsBackendRegistry {
	return &MemoryJsBackendRegistry{
		NewMemoryJsRegistry(),
		map[string]map[string]*Stamp{},
		expiDur,
	}
}

func (reg *MemoryJsBackendRegistry) StampedObject(dir, objName string, caStmp *Stamp) (*Object, *Stamp, error) {
	nameToStmp := reg.stmps[dir]
	if nameToStmp == nil {
		return nil, nil, nil
	}
	stmp := nameToStmp[objName]
	if stmp == nil {
		return nil, nil, nil
	}

	// 対象のスタンプを取得。

	newCaStmp := &Stamp{Date: stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: stmp.Digest}

	if caStmp != nil && !stmp.Date.After(caStmp.Date) && caStmp.Digest == stmp.Digest {
		return nil, newCaStmp, nil
	}

	// 無効なキャッシュだった。

	obj, _ := reg.Object(dir, objName)
	return obj, newCaStmp, nil
}
func (reg *MemoryJsBackendRegistry) AddObject(dir, objName string, obj *Object) error {
	reg.MemoryJsRegistry.AddObject(dir, objName, obj)
	nameToStmp := reg.stmps[dir]
	if nameToStmp == nil {
		nameToStmp = map[string]*Stamp{}
		reg.stmps[dir] = nameToStmp
	}
	nameToStmp[objName] = &Stamp{Date: time.Now(), Digest: strconv.Itoa(len(fmt.Sprint(obj)))}
	return nil
}
func (reg *MemoryJsBackendRegistry) RemoveObject(dir, objName string) error {
	reg.MemoryJsRegistry.RemoveObject(dir, objName)
	nameToStmp := reg.stmps[dir]
	if nameToStmp == nil {
		return nil
	}
	delete(nameToStmp, objName)
	return nil
}
