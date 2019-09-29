package generic_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/sets"
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

type Dog struct {
	ID      int
	Friends []string
}

func NewPerson() Dog {
	return Dog{
		ID:      3,
		Friends: sets.NewString().List(),
	}
}

func TestSomething(t *testing.T) {
	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")
}

func TestNewPerson(t *testing.T) {
	expected := Dog{
		ID:      3,
		Friends: sets.NewString().List(),
	}

	result := NewPerson()

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Unexpected mis-match: %s", diff.ObjectReflectDiff(expected, result))
	}
}
