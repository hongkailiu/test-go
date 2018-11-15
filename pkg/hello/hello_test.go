package main

import "testing"

func TestSum(t *testing.T) {
	err := NewHello()
	if err != nil {
		t.Errorf("error should not have occurred: %s", err.Error())
	}
}
