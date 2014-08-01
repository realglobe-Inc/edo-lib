package util

import (
	"fmt"
)

// なんちゃって木構造。

type Tree struct {
	isRoot func(label string) bool
	parent func(label string) string

	nodes map[string]*treeNode
}

func NewTree(isRoot func(string) bool, parent func(string) string) *Tree {
	return &Tree{
		isRoot,
		parent,
		map[string]*treeNode{},
	}
}

func (tree *Tree) String() string {
	return fmt.Sprint(tree.nodes)
}

type treeNode struct {
	val      interface{}
	children map[string]bool // 子のラベル。
}

func newTreeNode() *treeNode {
	return &treeNode{children: map[string]bool{}}
}

func (node *treeNode) String() string {
	return fmt.Sprint(node.val, toStringSlice(node.children))
}

func toStringSlice(set map[string]bool) []string {
	slice := []string{}
	for str, _ := range set {
		slice = append(slice, str)
	}
	return slice
}

// 通信形式から変換する。
func (tree *Tree) FromContainer(cont map[string]interface{}) {
	for label, val := range cont {
		tree.Add(label, val)
	}
}

// 通信形式へ変換する。
func (tree *Tree) ToContainer() (cont map[string]interface{}) {
	cont = map[string]interface{}{}
	for label, node := range tree.nodes {
		if node.val != nil {
			cont[label] = node.val
		}
	}
	return cont
}

// 使用形式の木にノードを加える。
// 親が居なければ勝手に生成する。
func (tree *Tree) Add(label string, val interface{}) {
	curNode := tree.nodes[label]
	if curNode == nil {
		curNode = newTreeNode()
		tree.nodes[label] = curNode
	}

	curNode.val = val

	// 親を作成。
	for curLabel := label; !tree.isRoot(curLabel); {
		parent := tree.parent(curLabel)
		parentNode := tree.nodes[parent]
		if parentNode == nil {
			parentNode = newTreeNode()
			tree.nodes[parent] = parentNode
		}
		parentNode.children[curLabel] = true
		curLabel = parent
	}
}

// 木からノードを削除する。
// ただし、子がいる場合は中間ノードとして残す。
func (tree *Tree) Remove(label string) {
	node := tree.nodes[label]
	if node == nil {
		return
	}
	for child, _ := range node.children {
		if tree.nodes[child] != nil {
			node.val = nil
			return
		}
	}
	delete(tree.nodes, label)

	// 親を削除。
	for curLabel := label; tree.isRoot(curLabel); {
		parent := tree.parent(curLabel)
		parentNode := tree.nodes[parent]
		delete(parentNode.children, curLabel)
		if parentNode.val != nil || len(parentNode.children) > 0 {
			break
		}
		curLabel = parent
	}
}

// 値を取り出す。
func (tree *Tree) Value(label string) (val interface{}) {
	node := tree.nodes[label]
	if node == nil {
		return nil
	}
	return node.val
}

// 使用形式の木から部分木 + αに含まれる値を列挙する。
func (tree *Tree) Values(label string) (vals []interface{}) {
	node := tree.nodes[label]
	if node == nil {
		return nil
	}

	vals = []interface{}{}
	tree._subTree(&vals, label)
	return vals
}

func (tree *Tree) _subTree(vals *[]interface{}, label string) {
	node := tree.nodes[label]
	if node == nil {
		return
	}

	if node.val != nil {
		*vals = append(*vals, node.val)
	}
	for child, _ := range node.children {
		tree._subTree(vals, child)
	}
}
