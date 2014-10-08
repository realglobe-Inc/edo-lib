package driver

import (
	"fmt"
	"github.com/realglobe-Inc/edo/util"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"reflect"
	"strings"
)

// イベントハンドラの管理。
type EventRegistry interface {
	// ハンドラを取得する。
	// イベントは / 区切りで木構造のノードを表し、そのノード以下の部分木に含まれる全てのハンドラを返す。
	Handler(usrUuid, event string, caStmp *Stamp) (hndl Handler, newCaStmp *Stamp, err error)

	// ハンドラを登録する。
	AddHandler(usrUuid, event string, hndl Handler) (newCaStmp *Stamp, err error)

	// ハンドラを削除する。
	RemoveHandler(usrUuid, event string) error
}

// 1 つのイベントに紐付けられるハンドラのリスト。
type Handler []*HandlerElement

// イベントハンドラ。
type HandlerElement struct {
	Url     string            `json:"url"               bson:"url"`
	Method  string            `json:"method,omitempty"  bson:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty" bson:"headers,omitempty"`
	Body    string            `json:"body,omitempty"    bson:"body,omitempty"`
	Rules   []*HandlerRule    `json:"rules,omitempty"   bson:"rules,omitempty"`
}

func (elem *HandlerElement) String() string {
	return fmt.Sprint(elem.Url+","+elem.Method+",", elem.Headers, ",", len(elem.Body), ",", len(elem.Rules))
}

// イベントの付属パラメータの扱いを記述する予定。
type HandlerRule struct {
}

// 骨組み。
// バックエンドでは、ユーザーごとにイベントからハンドラへの写像を保存。
// それを木構造に構成し直して使う。
type eventRegistry struct {
	base KeyValueStore
}

func newEventRegistry(base KeyValueStore) *eventRegistry {
	return &eventRegistry{base}
}

func (reg *eventRegistry) Handler(usrUuid, event string, stmp *Stamp) (hndl Handler, newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrUuid, stmp)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	} else if value == nil {
		return nil, newCaStmp, nil
	}

	tree := newEventTree()
	tree.fromContainer(value.(map[string]Handler))

	return tree.handler(event), newCaStmp, nil
}

func (reg *eventRegistry) AddHandler(usrUuid, event string, hndl Handler) (newCaStmp *Stamp, err error) {
	value, newCaStmp, err := reg.base.Get(usrUuid, nil)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	var eventToHndl map[string]Handler
	if value == nil {
		eventToHndl = map[string]Handler{event: hndl}
	} else {
		eventToHndl := value.(map[string]Handler)
		if reflect.DeepEqual(eventToHndl[event], hndl) {
			// 書き込みは減らす。
			return newCaStmp, nil
		}
	}

	return reg.base.Put(usrUuid, eventToHndl)
}

func (reg *eventRegistry) RemoveHandler(usrUuid, event string) error {
	value, _, err := reg.base.Get(usrUuid, nil)
	if err != nil {
		return erro.Wrap(err)
	} else if value == nil {
		return nil
	}

	eventToHndl := value.(map[string]Handler)
	if _, ok := eventToHndl[event]; !ok {
		// 書き込みは減らす。
		return nil
	}
	delete(eventToHndl, event)

	if len(eventToHndl) > 0 {
		if _, err := reg.base.Put(usrUuid, eventToHndl); err != nil {
			return erro.Wrap(err)
		}
	} else if err := reg.base.Remove(usrUuid); err != nil {
		return erro.Wrap(err)
	}
	return nil
}

// EventRegistry 用内部データ。
type eventTree struct {
	*util.Tree
}

func newEventTree() *eventTree {
	return &eventTree{util.NewTree(
		func(label string) bool {
			return label == "/"
		},
		func(label string) string {
			if idx := strings.LastIndex(label, "/"); idx == 0 {
				return "/"
			} else {
				return label[:idx]
			}
		},
	)}
}

func (tree *eventTree) add(event string, hndl Handler) {
	tree.Add(event, hndl)
}

func (tree *eventTree) remove(event string) {
	tree.Remove(event)
}

func (tree *eventTree) handler(event string) (hndl Handler) {
	values := tree.Values(event)
	if values == nil {
		return nil
	}
	hndl = Handler{}
	for _, value := range values {
		for _, e := range value.(Handler) {
			hndl = append(hndl, e)
		}
	}
	return hndl
}

func (tree *eventTree) fromContainer(cont map[string]Handler) {
	c := map[string]interface{}{}
	for event, hndl := range cont {
		c[event] = hndl
	}
	tree.FromContainer(c)
}

func (tree *eventTree) toContainer() (cont map[string]Handler) {
	c := tree.ToContainer()
	cont = map[string]Handler{}
	for label, value := range c {
		cont[label] = value.(Handler)
	}
	return cont
}
