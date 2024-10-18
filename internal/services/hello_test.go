package services

import (
	"testing"
)

func TestGetHelloMessage(t *testing.T) {
	expected := "Hello, World!"
	actual := GetHelloMessage()

	if actual != expected {
		t.Errorf("GetHelloMessage() = %v; want %v", actual, expected)
	}
}
