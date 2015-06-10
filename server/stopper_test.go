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
	"sync/atomic"
	"testing"
	"time"
)

func TestStopper(t *testing.T) {
	loop := 100
	proc := 100
	for i := 0; i < loop; i++ {
		var n int64 = 0

		s := NewStopper()
		for j := 0; j < proc; j++ {
			s.Stop()
			go func() {
				time.Sleep(time.Millisecond)
				atomic.AddInt64(&n, 1)
				s.Unstop()
			}()
		}

		s.Lock()
		defer s.Unlock()
		for s.Stopped() {
			s.Wait()
		}

		if n != int64(proc) {
			t.Error(n)
			t.Fatal(proc)
		}
	}
}
