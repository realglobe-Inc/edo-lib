package driver

import (
	"fmt"
	"strconv"
	"time"
)

// JavaScript.
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

// ID プロバイダ。
type MemoryIdProviderBackend struct {
	*MemoryIdProviderRegistry
	stmp    *Stamp
	expiDur time.Duration
}

func NewMemoryIdProviderBackend(expiDur time.Duration) *MemoryIdProviderBackend {
	stmp := &Stamp{Date: time.Now(), Digest: strconv.Itoa(0)}
	stmp.ExpiDate = stmp.Date.Add(expiDur)

	return &MemoryIdProviderBackend{NewMemoryIdProviderRegistry(), stmp, expiDur}
}

func (reg *MemoryIdProviderBackend) StampedIdProviders(caStmp *Stamp) ([]*IdProvider, *Stamp, error) {
	newCaStmp := &Stamp{Date: reg.stmp.Date, ExpiDate: time.Now().Add(reg.expiDur), Digest: reg.stmp.Digest}

	if caStmp == nil || caStmp.Date.Before(reg.stmp.Date) || caStmp.Digest != reg.stmp.Digest {
		idps, _ := reg.IdProviders()
		return idps, newCaStmp, nil
	}

	return nil, newCaStmp, nil
}
func (reg *MemoryIdProviderBackend) AddIdProvider(idp *IdProvider) {
	reg.MemoryIdProviderRegistry.AddIdProvider(idp)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
func (reg *MemoryIdProviderBackend) RemoveIdProvider(idpUuid string) {
	reg.MemoryIdProviderRegistry.RemoveIdProvider(idpUuid)
	dig, _ := strconv.Atoi(reg.stmp.Digest)
	reg.stmp = &Stamp{Date: time.Now(), Digest: strconv.Itoa(dig + 1)}
}
