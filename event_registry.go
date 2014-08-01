package driver

import (
	"fmt"
	"github.com/realglobe-Inc/edo/util"
	"strings"
)

type EventRegistry interface {
	// ハンドラを取得する。
	// イベントは / 区切りで木構造のノードを表し、そのノード以下の部分木に含まれる全てのハンドラを返す。
	Handler(usrUuid, event string) (Handler, error)

	// ハンドラを登録する。
	AddHandler(usrUuid, event string, hndl Handler) error
	// ハンドラを削除する。
	RemoveHandler(usrUuid, event string) error
}

type Handler []*HandlerElement

type HandlerElement struct {
	Url     string            `json:"url"               bson:"url"`
	Method  string            `json:"method,omitempty"  bson:"method,omitempty"`
	Headers map[string]string `json:"headers,omitempty" bson:"headers,omitempty"`
	Body    string            `json:"body,omitempty"    bson:"body,omitempty"`
	Rules   []*HandlerRule    `json:"rules,omitempty"   bson:"rules,omitempty"`
}

func (elem *HandlerElement) String() string {
	return fmt.Sprint("{"+elem.Url+" "+elem.Method+" ", elem.Headers, " ", len(elem.Body), " ", len(elem.Rules), "}")
}

// イベントの付属パラメータの扱いを記述する予定。
type HandlerRule struct {
}

// EventRegistry の内部データ。

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
	vals := tree.Values(event)
	if vals == nil {
		return nil
	}
	hndl = Handler{}
	for _, val := range vals {
		for _, e := range val.(Handler) {
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
	for label, val := range c {
		cont[label] = val.(Handler)
	}
	return cont
}
