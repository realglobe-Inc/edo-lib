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

package server

import (
	"sync"
)

// 処理途中でサーバーが終了するのを防ぐ。
//
// メイン側で、
//  s.Lock()
//  defer s.Unlock()
//  for s.Stopped() {
//  	s.Wait()
//  }
//
// 処理側で、
//  s.Stop()
//  defer s.Unstop()
type Stopper struct {
	c *sync.Cond
	n int
}

func NewStopper() *Stopper {
	return &Stopper{
		c: sync.NewCond(&sync.Mutex{}),
	}
}

func (this *Stopper) Lock() {
	this.c.L.Lock()
}

func (this *Stopper) Unlock() {
	this.c.L.Unlock()
}

func (this *Stopper) Stopped() bool {
	return this.n > 0
}

func (this *Stopper) Wait() {
	this.c.Wait()
}

func (this *Stopper) Stop() {
	this.c.L.Lock()
	defer this.c.L.Unlock()
	this.n++
}

func (this *Stopper) Unstop() {
	this.c.L.Lock()
	this.n--
	n := this.n
	this.c.L.Unlock()
	if n == 0 {
		this.c.Signal()
	}
}
