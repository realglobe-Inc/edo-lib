// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// キャッシュ。
package cache

import (
	"container/heap"
)

// キャッシュ。
type Cache interface {
	// 入れる。
	Put(key, val, prio interface{})
	// 取り出す。
	Get(key interface{}) (val, prio interface{})
	// 優先度を変えつつ取り出す。LRU のとき使う。
	Update(key, prio interface{}) (val interface{})
	// 基準以下を削除。
	// 優先度 nil はいかなる非 nil な優先度より低いとする。
	// よって、Update で優先度を nil にしてから、CleanLower すれば削除できる。
	// nil で CleanLower したときは優先度 nil のものだけを削除する。
	CleanLower(prioThres interface{})
}

// スレッドセーフではないキャッシュを返す。
// less は非 nil の優先度 2 つを比べる関数。優先度 nil に対する挙動は指定できない。
func New(less func(interface{}, interface{}) bool) Cache {
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
	idx, ok := ca.keyToIdx[key]
	if !ok {
		heap.Push(ca, &cacheElement{key, val, prio})
		return
	}
	elem := ca.prioQueue[idx]
	elem.val = val
	elem.prio = prio
	heap.Fix(ca, idx)
	return
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
	key  interface{}
	val  interface{}
	prio interface{}
}

func (ca *cache) Len() int {
	return len(ca.prioQueue)
}

func (ca *cache) Less(i, j int) bool {
	prio1, prio2 := ca.prioQueue[i].prio, ca.prioQueue[j].prio
	if prio1 == nil {
		// nil < 非 nil。
		// nil == nil。
		return prio2 != nil
	} else if prio2 == nil {
		// 非 nil > nil。
		return false
	} else {
		// 非 nil ? 非 nil。
		return ca.less(prio1, prio2)
	}
}
func (ca *cache) Swap(i, j int) {
	elem1, elem2 := ca.prioQueue[i], ca.prioQueue[j]
	ca.prioQueue[i], ca.prioQueue[j] = elem2, elem1
	ca.keyToIdx[elem1.key], ca.keyToIdx[elem2.key] = j, i
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
