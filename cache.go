package util

import (
	"container/heap"
)

// キャッシュ。
// スレッドセーフではない。

type Cache interface {
	// 入れる。
	Put(key, val, prio interface{})
	// 取り出す。
	Get(key interface{}) (val, prio interface{})
	// 優先度を変えつつ取り出す。LRU のとき使う。
	Update(key, prio interface{}) (val interface{})
	// 基準以下を削除。
	CleanLesser(prioThres interface{})
}

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

func (ca *cache) Put(key, val, prio interface{}) {
	ca.Push(&cacheElement{key, val, prio})
}

func (ca *cache) Get(key interface{}) (val, prio interface{}) {
	idx, ok := ca.keyToIdx[key]
	if !ok {
		return nil, nil
	}
	elem := ca.prioQueue[idx]
	return elem.val, elem.prio
}

func (ca *cache) Update(key, prio interface{}) (val interface{}) {
	idx, ok := ca.keyToIdx[key]
	if !ok {
		return nil
	}
	elem := ca.prioQueue[idx]
	elem.prio = prio
	heap.Fix(ca, idx)
	return elem.val
}

func (ca *cache) CleanLesser(prioThres interface{}) {
	for ca.Len() > 0 {
		if ca.less(ca.prioQueue[0].prio, prioThres) {
			heap.Pop(ca)
		} else {
			break
		}
	}
}

type cacheElement struct {
	key  interface{}
	val  interface{}
	prio interface{}
}

func (ca *cache) Len() int {
	return len(ca.prioQueue)
}

func (ca *cache) Less(i, j int) bool {
	return ca.less(ca.prioQueue[i].prio, ca.prioQueue[j].prio)
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
