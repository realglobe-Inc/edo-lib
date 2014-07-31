package driver

import (
	"fmt"
	"strconv"
	"time"
)

// JavaScript.
type MemoryJsBackendRegistry struct {
	*MemoryJsRegistry
	stmps map[string]map[string]*Stamp
}

func NewMemoryJsBackendRegistry() *MemoryJsBackendRegistry {
	return &MemoryJsBackendRegistry{
		NewMemoryJsRegistry(),
		map[string]map[string]*Stamp{},
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

	newCaStmp := &Stamp{Date: time.Now(), Digest: stmp.Digest}

	if caStmp != nil && caStmp.Date.After(stmp.Date) && caStmp.Digest == stmp.Digest {
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
