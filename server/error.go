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
	"net/http"
	"strconv"
)

// HTTP のステータスコードを付加したエラー。
type Error struct {
	status int
	msg    string

	cause error
}

// stat が 0 の場合、代わりに http.StatusInternalServerError が入る。
func NewError(stat int, msg string, cause error) *Error {
	if stat <= 0 {
		stat = http.StatusInternalServerError
	}
	return &Error{stat, msg, cause}
}

func (this *Error) Error() string {
	prefix := ""
	if this.cause != nil {
		prefix += this.cause.Error() + "\ncaused "
	}
	return prefix + strconv.Itoa(this.status) + " " + http.StatusText(this.status) + ": " + this.msg
}

func (this *Error) Status() int {
	return this.status
}

func (this *Error) Message() string {
	return this.msg
}

func (this *Error) Cause() error {
	return this.cause
}
