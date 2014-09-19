package util

import (
	"encoding/json"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"time"
)

// "10s" みたいな形式で JSON にできる time.Duration のラッパー。

type Duration struct {
	time.Duration
}

func (dur *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(dur.Duration.String())
}

func (dur *Duration) UnmarshalJSON(b []byte) error {
	var buff string
	if err := json.Unmarshal(b, &buff); err != nil {
		return erro.Wrap(err)
	}
	var err error
	dur.Duration, err = time.ParseDuration(buff)
	if err != nil {
		return erro.Wrap(err)
	}
	return nil
}

func NewDuration(dur time.Duration) *Duration {
	return &Duration{dur}
}

func ParseDuration(s string) (*Duration, error) {
	dur, err := time.ParseDuration(s)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return &Duration{dur}, nil
}
