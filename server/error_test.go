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
	"errors"
	"net/http"
	"testing"
)

func TestStatusError(t *testing.T) {
	cause := errors.New("test error")
	err := NewStatusError(http.StatusBadRequest, "test error", cause)

	if err.Status() != http.StatusBadRequest {
		t.Fatal(err.Status(), http.StatusBadRequest)
	} else if err.Message() != "test error" {
		t.Fatal(err.Message(), "test error")
	} else if err.Cause() != cause {
		t.Fatal(err.Cause(), cause)
	} else if len(err.Error()) <= len(err.Message())+len(cause.Error()) {
		t.Error(err.Error())
		t.Error(err.Message())
		t.Fatal(cause.Error())
	}
}
