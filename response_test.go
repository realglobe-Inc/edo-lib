package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestErrorToResponseJson(t *testing.T) {
	before := errors.New("abcde")
	errJson := ErrorToResponseJson(before)

	var after struct {
		Status  int
		Message string
	}
	if err := json.Unmarshal(errJson, &after); err != nil {
		t.Fatal(err, string(errJson))
	}
	if int(after.Status) != http.StatusInternalServerError {
		t.Error(after)
	} else if after.Message != "abcde" {
		t.Error(after)
	}
}

func TestHttpStatusErrorToResponseJson(t *testing.T) {
	before := NewHttpStatusError(http.StatusNotFound, "abcde", nil)
	errJson := ErrorToResponseJson(before)

	var after struct {
		Status  int
		Message string
	}
	if err := json.Unmarshal(errJson, &after); err != nil {
		t.Fatal(err, string(errJson))
	}
	if after.Status != http.StatusNotFound {
		t.Error(after)
	} else if after.Message != "abcde" {
		t.Error(after)
	}
}
