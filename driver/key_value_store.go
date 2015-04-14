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

package driver

import (
	"io"
)

type KeyValueStore interface {
	Get(key string, caStmp *Stamp) (val interface{}, newCaStmp *Stamp, err error)
	Put(key string, val interface{}) (*Stamp, error)
	Remove(key string) error

	io.Closer
}