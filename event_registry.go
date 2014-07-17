package driver

import (
	"fmt"
	"strings"
)

// EventRegistry の内部データ。

type eventTree map[string]*eventNode

type eventNode struct {
	hndl        Handler
	childEvents map[string]bool
}

func (node *eventNode) String() string {
	return fmt.Sprint(node.hndl, toStringSlice(node.childEvents))
}

func newEventNode() *eventNode {
	return &eventNode{childEvents: map[string]bool{}}
}

// 通信形式から変換する。
func (tree eventTree) fromContainer(cont map[string]Handler) {
	for event, hndl := range cont {
		tree.add(event, hndl)
	}
}

// 通信形式へ変換する。
func (tree eventTree) toContainer() (cont map[string]Handler) {
	cont = map[string]Handler{}
	for event, node := range tree {
		if node.hndl != nil {
			cont[event] = node.hndl
		}
	}
	return cont
}

// 使用形式の木にノードを加える。
// 親が居なければ勝手に生成する。
func (tree eventTree) add(event string, hndl Handler) {
	curNode := tree[event]
	if curNode == nil {
		curNode = newEventNode()
		tree[event] = curNode
	}

	curNode.hndl = hndl

	// 親を作成。
	for curEvent := event; curEvent != "/"; {
		var parentEvent string
		if idx := strings.LastIndex(curEvent, "/"); idx == 0 {
			parentEvent = "/"
		} else {
			parentEvent = curEvent[:idx]
		}
		parent := tree[parentEvent]
		if parent == nil {
			parent = newEventNode()
			tree[parentEvent] = parent
		}
		parent.childEvents[curEvent] = true
		curEvent = parentEvent
	}
}

// 木からノードを削除する。
// ただし、子がいる場合は中間ノードとして残す。
func (tree eventTree) remove(event string) {
	node := tree[event]
	if node == nil {
		return
	}
	for childEvent, _ := range node.childEvents {
		if tree[childEvent] != nil {
			node.hndl = nil
			return
		}
	}
	delete(tree, event)

	// 親を削除。
	for curEvent := event; curEvent != "/"; {
		var parentEvent string
		if idx := strings.LastIndex(curEvent, "/"); idx == 0 {
			parentEvent = "/"
		} else {
			parentEvent = curEvent[:idx]
		}
		parent := tree[parentEvent]
		delete(parent.childEvents, curEvent)
		if parent.hndl != nil || len(parent.childEvents) > 0 {
			break
		}
		curEvent = parentEvent
	}
}

// 使用形式の木から部分木 + αに含まれるアドレスを列挙する。
func (tree eventTree) handler(event string) (hndl Handler) {
	node := tree[event]
	if node == nil {
		return nil
	}

	hndl = Handler{}
	tree._subTree(&hndl, event)
	return hndl
}

func (tree eventTree) _subTree(hndl *Handler, event string) {
	node := tree[event]
	if node == nil {
		return
	}

	if node.hndl != nil {
		*hndl = append(*hndl, node.hndl...)
	}
	for childEvent, _ := range node.childEvents {
		tree._subTree(hndl, childEvent)
	}
}
