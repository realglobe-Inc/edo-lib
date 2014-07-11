package util

import (
	"testing"
)

func TestHasApiPath(t *testing.T) {
	if HasApiPath("/", "/") {
		t.Error()
	}
	if !HasApiPath("/a", "/") {
		t.Error()
	}
	if HasApiPath("/", "/a") {
		t.Error()
	}
	if HasApiPath("/a", "/a") {
		t.Error()
	}
	if !HasApiPath("/a/b", "/a") {
		t.Error()
	}
}

func TestTrimApiPath(t *testing.T) {
	if b := TrimApiPath("/a", "/"); b != "/a" {
		t.Error(b)
	}
	if b := TrimApiPath("/a/b", "/a"); b != "/b" {
		t.Error(b)
	}
}
