package lc

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFindMedianSortedArrays1(t *testing.T) {
	o := NewGomegaWithT(t)

	nums1 := []int{1, 3}
	nums2 := []int{2}
	r := findMedianSortedArrays(nums1, nums2)

	o.Expect(r).To(Equal(2.0))
}

func TestFindMedianSortedArrays2(t *testing.T) {
	o := NewGomegaWithT(t)

	nums1 := []int{1, 2}
	nums2 := []int{3, 4}
	r := findMedianSortedArrays(nums1, nums2)

	o.Expect(r).To(Equal(2.5))
}
