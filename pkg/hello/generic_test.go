package main_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person interface {
	getName() string
}

type Student struct {
	name string
}

type Freshman struct {
	Student
}

func (s *Student) getName() string {
	return "a: " + s.name
}

func (f *Freshman) getName() string {
	return "b: " + f.name
}

func Test1(t *testing.T) {
	var persons []Person
	s := Student{name: "hongkai"}
	persons = append(persons, &s)
	f := Freshman{Student{name: "elle"}}
	persons = append(persons, &f)

	for _, p := range persons {
		fmt.Println(p.getName())
	}
}

func TestSomething(t *testing.T) {
	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")
}
