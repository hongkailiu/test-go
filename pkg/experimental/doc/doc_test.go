package doc_test

import (
	"testing"
)

func TestNewHello(t *testing.T) {
	err := NewHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}

func TestNewNiceHello(t *testing.T) {
	err := MYNewNiceHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}
