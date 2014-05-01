package main

import (
	"testing"
)

func AssertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if actual != expected {
		t.Errorf("Expected '%s' but received '%s'", expected, actual)
	}
}

func AssertNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if actual == expected {
		t.Errorf("Expected '%s' to be different from '%s'", expected, actual)
	}
}

func AssertNotError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error not expected: %s", err.Error())
	}
}
