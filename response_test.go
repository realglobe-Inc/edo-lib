package util

import (
	"encoding/json"
	"errors"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"testing"
)

func TestErrorToResponseJson(t *testing.T) {
	before := errors.New("aaa")
	errJson := ErrorToResponseJson(before)

	var after map[string]interface{}
	if err := json.Unmarshal(errJson, &after); err != nil {
		t.Fatal(err, string(errJson))
	}
	if after["name"] != "Error" {
		t.Error(after)
	} else if after["message"] != "aaa" {
		t.Error(after)
	} else if after["sys_type"] == "" {
		t.Error(after)
	} else if _, ok := after["sys_data"]; !ok {
		t.Error(after)
	} else if _, ok := after["sys_stack"]; ok {
		t.Error(after)
	}
}

func TestErrorToResponseJsonWithStack(t *testing.T) {
	before := erro.New("aaa")
	errJson := ErrorToResponseJson(before)

	var after map[string]interface{}
	if err := json.Unmarshal(errJson, &after); err != nil {
		t.Fatal(err, string(errJson))
	}
	if after["name"] != "Error" {
		t.Error(after)
	} else if after["message"] != "aaa" {
		t.Error(after)
	} else if after["sys_type"] == "" {
		t.Error(after)
	} else if _, ok := after["sys_data"]; !ok {
		t.Error(after)
	} else if _, ok := after["sys_stack"]; !ok {
		t.Error(after)
	}
}
