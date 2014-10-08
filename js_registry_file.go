package driver

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
	"time"
)

// value を string として返す。
func stringMarshal(value interface{}) ([]byte, error) {
	return []byte(value.(string)), nil
}

// data を string として返す。
func stringUnmarshal(data []byte) (interface{}, error) {
	return string(data), nil
}

// コード本体と補助情報を分けて保存する。
type objectHeader struct {
	Service bool     `json:"service,omitempty"`
	Library bool     `json:"library,omitempty"`
	Include []string `json:"include,omitempty"`
}

// data を JSON として、objectHeader にデコードする。
func objectHeaderUnmarshal(data []byte) (interface{}, error) {
	var res objectHeader
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, erro.Wrap(err)
	}
	return &res, nil
}

type fileJsRegistry struct {
	code   KeyValueStore
	header KeyValueStore
}

// スレッドセーフ。
func NewFileJsRegistry(path string, expiDur time.Duration) JsRegistry {
	return newSynchronizedFileJsRegistry(newFileJsRegistry(path, expiDur))
}

func jsKeyGen(before string) string {
	return before + ".js"
}

// スレッドセーフではない。
func newFileJsRegistry(path string, expiDur time.Duration) *fileJsRegistry {
	return &fileJsRegistry{
		newCachingKeyValueStore(newFileKeyValueStore(path, jsKeyGen, stringMarshal, stringUnmarshal, expiDur)),
		newCachingKeyValueStore(newFileKeyValueStore(path, jsonKeyGen, json.Marshal, objectHeaderUnmarshal, expiDur)),
	}
}

func (reg *fileJsRegistry) Object(dir, objName string, caStmp *Stamp) (obj *Object, newCaStmp *Stamp, err error) {
	var code string
	if value, stmp, err := reg.code.Get(dir+"/"+objName, caStmp); err != nil {
		return nil, nil, erro.Wrap(err)
	} else if stmp == nil {
		return nil, nil, nil
	} else {
		if value != nil {
			code = value.(string)
		}
		newCaStmp = stmp
	}

	if value, stmp, err := reg.header.Get(dir+"/"+objName, nil); err != nil {
		return nil, nil, erro.Wrap(err)
	} else {
		if value == nil {
			obj = &Object{Code: code}
		} else {
			header := value.(*objectHeader)
			obj = &Object{header.Service, header.Library, header.Include, code}
		}
		newCaStmp.Digest += stmp.Digest
	}

	if caStmp != nil && !newCaStmp.Date.After(caStmp.Date) && caStmp.Digest == newCaStmp.Digest {
		return nil, newCaStmp, nil
	}

	return obj, newCaStmp, nil
}

func (reg *fileJsRegistry) AddObject(dir, objName string, obj *Object) (*Stamp, error) {
	newCaStmp, err := reg.code.Put(dir+"/"+objName, obj.Code)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	if obj.Service || obj.Library || len(obj.Include) > 0 {
		if _, err := reg.header.Put(dir+"/"+objName, &objectHeader{obj.Service, obj.Library, obj.Include}); err != nil {
			return nil, erro.Wrap(err)
		}
	}
	return newCaStmp, nil
}

func (reg *fileJsRegistry) RemoveObject(dir, objName string) error {
	if err := reg.code.Remove(dir + "/" + objName); err != nil {
		return erro.Wrap(err)
	}
	if err := reg.header.Remove(dir + "/" + objName); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// スレッドセーフ化。
type synchronizedJsRegistry synchronizedDriver

type objectRequest struct {
	dir     string
	objName string
	caStmp  *Stamp

	objCh       chan *Object
	newCaStmpCh chan *Stamp
}

type addObjectRequest struct {
	dir     string
	objName string
	obj     *Object

	newCaStmpCh chan *Stamp
}

type removeObjectRequest struct {
	dir     string
	objName string
}

func newSynchronizedFileJsRegistry(base JsRegistry) JsRegistry {
	return (*synchronizedJsRegistry)(newSynchronizedDriver(map[reflect.Type]func(interface{}, chan<- error){
		reflect.TypeOf(&objectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*objectRequest)
			obj, stmp, err := base.Object(req.dir, req.objName, req.caStmp)
			if err != nil {
				errCh <- err
			} else {
				req.objCh <- obj
				req.newCaStmpCh <- stmp
			}
		},
		reflect.TypeOf(&addObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*addObjectRequest)
			stmp, err := base.AddObject(req.dir, req.objName, req.obj)
			if err != nil {
				errCh <- err
			} else {
				req.newCaStmpCh <- stmp
			}
		},
		reflect.TypeOf(&removeObjectRequest{}): func(r interface{}, errCh chan<- error) {
			req := r.(*removeObjectRequest)
			errCh <- base.RemoveObject(req.dir, req.objName)
		},
	}))
}

func (reg *synchronizedJsRegistry) Object(dir, objName string, caStmp *Stamp) (obj *Object, newCaStmp *Stamp, err error) {
	objCh := make(chan *Object, 1)
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&objectRequest{dir, objName, caStmp, objCh, newCaStmpCh}, errCh}
	select {
	case obj := <-objCh:
		return obj, <-newCaStmpCh, nil
	case err := <-errCh:
		return nil, nil, err
	}
}

func (reg *synchronizedJsRegistry) AddObject(dir, objName string, obj *Object) (newCaStmp *Stamp, err error) {
	newCaStmpCh := make(chan *Stamp, 1)
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&addObjectRequest{dir, objName, obj, newCaStmpCh}, errCh}
	select {
	case stmp := <-newCaStmpCh:
		return stmp, nil
	case err := <-errCh:
		return nil, err
	}
}

func (reg *synchronizedJsRegistry) RemoveObject(dir, objName string) error {
	errCh := make(chan error, 1)
	reg.reqCh <- &synchronizedRequest{&removeObjectRequest{dir, objName}, errCh}
	return <-errCh
}
