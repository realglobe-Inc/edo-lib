package util

import (
	"container/heap"
)

// キャッシュ。
// スレッドセーフではない。

type Cache interface {
	// 入れる。
	Put(key, value, prio interface{})
	// 取り出す。
	Get(key interface{}) (value, prio interface{})
	// 優先度を変えつつ取り出す。LRU のとき使う。
	Update(key, prio interface{}) (value interface{})
	// 基準以下を削除。
	// 優先度 nil はいかなる非 nil な優先度より低いとする。
	// よって、Update で優先度を nil にしてから、CleanLower すれば削除できる。
	// nil で CleanLower したときは優先度を nil のものだけを削除する。
	CleanLower(prioThres interface{})
}

// less は非 nil の優先度 2 つを比べる関数。優先度 nil に対する挙動は指定できない。
func NewCache(less func(interface{}, interface{}) bool) Cache {
	ca := &cache{less, []*cacheElement{}, map[interface{}]int{}}
	heap.Init(ca)
	return ca
}

type cache struct {
	// Cache と heap.Interface を実装する。
	// Cache からは参照だけ直接行い、操作は heap.Interface 側でやる。

	less      func(interface{}, interface{}) bool
	prioQueue []*cacheElement
	keyToIdx  map[interface{}]int
}

func (ca *cache) Put(key, value, prio interface{}) {
	idx, ok := ca.keyToIdx[key]
	if !ok {
		ca.Push(&cacheElement{key, value, prio})
		return
	}
	elem := ca.prioQueue[idx]
	elem.value = value
	elem.prio = prio
	heap.Fix(ca, idx)
	return
}

func (ca *cache) Get(key interface{}) (value, prio interface{}) {
	idx, ok := ca.keyToIdx[key]
	if !ok {
		return nil, nil
	}
	elem := ca.prioQueue[idx]
	return elem.value, elem.prio
}

func (ca *cache) Update(key, prio interface{}) (value interface{}) {
	idx, ok := ca.keyToIdx[key]
	if !ok {
		return nil
	}
	elem := ca.prioQueue[idx]
	elem.prio = prio
	heap.Fix(ca, idx)
	return elem.value
}

func (ca *cache) CleanLower(prioThres interface{}) {
	for ca.Len() > 0 {
		if ca.prioQueue[0].prio == nil {
			heap.Pop(ca)
		} else if prioThres == nil {
			break
		} else if ca.less(ca.prioQueue[0].prio, prioThres) {
			heap.Pop(ca)
		} else {
			break
		}
	}
}

type cacheElement struct {
	key   interface{}
	value interface{}
	prio  interface{}
}

func (ca *cache) Len() int {
	return len(ca.prioQueue)
}

func (ca *cache) Less(i, j int) bool {
	if ca.prioQueue[i].prio == nil {
		// nil < 非 nil。
		// nil == nil。
		return ca.prioQueue[j].prio != nil
	} else if ca.prioQueue[j].prio == nil {
		// 非 nil > nil。
		return false
	} else {
		// 非 nil ? 非 nil。
		return ca.less(ca.prioQueue[i].prio, ca.prioQueue[j].prio)
	}
}
func (ca *cache) Swap(i, j int) {
	ca.prioQueue[i], ca.prioQueue[j] = ca.prioQueue[j], ca.prioQueue[i]
	ca.keyToIdx[ca.prioQueue[i].key], ca.keyToIdx[ca.prioQueue[j].key] = i, j
}

func (ca *cache) Push(x interface{}) {
	elem := x.(*cacheElement)
	ca.prioQueue = append(ca.prioQueue, elem)
	ca.keyToIdx[elem.key] = len(ca.prioQueue) - 1
}

func (ca *cache) Pop() interface{} {
	n := len(ca.prioQueue)
	elem := ca.prioQueue[n-1]
	ca.prioQueue = ca.prioQueue[:n-1]
	delete(ca.keyToIdx, elem.key)
	return elem
}
