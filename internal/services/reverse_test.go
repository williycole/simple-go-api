package services

import (
	"testing"
)

func TestReverseMessage(t *testing.T) {
	s := "Hello, World!"
	expected := "!dlroW ,olleH"
	actual := GetReverseMessage(s)

	if actual != expected {
		t.Errorf("GetReverseMessage() = %v; want %v", actual, expected)
	}

}
