package lc

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test1(t *testing.T) {
	o := NewGomegaWithT(t)

	r := twoSum(nil, 0)

	o.Expect(r).To(BeNil())
}

func Test2(t *testing.T) {
	o := NewGomegaWithT(t)

	r := twoSum([]int{2, 7, 11, 15}, 9)

	o.Expect(r).To(Equal([]int{0, 1}))
}


func Test3(t *testing.T) {
	o := NewGomegaWithT(t)

	r := twoSum([]int{3, 2, 4}, 6)

	o.Expect(r).To(Equal([]int{1, 2}))
}