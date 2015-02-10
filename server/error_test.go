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
		t.Error(err.Status(), http.StatusBadRequest)
	} else if err.Message() != "test error" {
		t.Error(err.Message(), "test error")
	} else if err.Cause() != cause {
		t.Error(err.Cause(), cause)
	} else if len(err.Error()) <= len(err.Message())+len(cause.Error()) {
		t.Error(err.Error())
		t.Error(err.Message())
		t.Error(cause.Error())
	}
}
