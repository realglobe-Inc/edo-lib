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
	"fmt"
)

// HTTP のステータスコードを付加したエラー。
type StatusError struct {
	status int
	msg    string

	cause error
}

func NewStatusError(status int, msg string, cause error) *StatusError {
	return &StatusError{status, msg, cause}
}

func (err *StatusError) Error() string {
	buff := err.msg
	if err.cause != nil {
		buff += fmt.Sprintln()
		buff += "caused by: "
		buff += err.cause.Error()
	}
	return buff
}

func (err *StatusError) Status() int {
	return err.status
}

func (err *StatusError) Message() string {
	return err.msg
}

func (err *StatusError) Cause() error {
	return err.cause
}
