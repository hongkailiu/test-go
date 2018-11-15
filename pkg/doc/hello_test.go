package doc_test

import (
	"testing"

	. "github.com/hongkailiu/test-go/pkg/doc"
)

func TestNewHello(t *testing.T) {
	err := NewHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}

func TestNewNiceHello(t *testing.T) {
	err := NewNiceHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}
