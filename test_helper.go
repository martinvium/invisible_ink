package main

import (
	"testing"
)

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, actual)
	}
}
