package main

import (
	"testing"

	"golang.org/x/example/stringutil"
)

func TestReverse(t *testing.T) {
	if stringutil.Reverse("Hello, world") != "dlrow ,olleH" {
		t.Error("Error reverse")
	}
}
