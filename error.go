package util

import (
	"fmt"
)

type PanicWrapper struct {
	Recv  interface{} `json:"recover"` // recover の返り値。
	stack string
}

func (err *PanicWrapper) Error() string {
	return fmt.Sprint(err.Recv)
}

func NewPanicWrapper(recv interface{}, stack string) *PanicWrapper {
	return &PanicWrapper{recv, stack}
}
