package lc

import (
	"fmt"
)

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	var tree *node
	tree = insert(nums1, nums2, tree)

	if d(tree.l) == d(tree.r) {
		return float64(tree.v)
	}
	return float64((float64(tree.v + tree.l.v)) / 2)
}

func insert(nums1 []int, nums2 []int, tree *node) *node {

	if len(nums1) == 0 && len(nums2) == 0 {
		return tree
	}
	var v int
	var flag bool
	if len(nums1) != 0 && len(nums2) == 0 {
		v = nums1[0]
		flag = true
	}
	if len(nums1) == 0 && len(nums2) != 0 {
		v = nums2[0]
	}
	if len(nums1) != 0 && len(nums2) != 0 {
		v = nums1[0]
		flag = true
		if nums2[0] < nums1[0] {
			v = nums2[0]
			flag = false
		}
	}
	newNode := &node{v, nil, nil}
	//fmt.Printf("chose v: %d\n", v)

	if tree == nil {
		tree = newNode
		if flag {
			return insert(nums1[1:], nums2, tree)
		}
		return insert(nums1, nums2[1:], tree)
	}

	if tree.l == nil {
		newNode.l = tree
		tree = newNode
		if flag {
			return insert(nums1[1:], nums2, tree)
		}
		return insert(nums1, nums2[1:], tree)
	}

	p := tree
	for p.r != nil {
		p = p.r
	}
	p.r = newNode

	if d(tree.l) == d(tree.r) {
		if flag {
			return insert(nums1[1:], nums2, tree)
		}
		return insert(nums1, nums2[1:], tree)
	}

	if d(tree.l) < d(tree.r) {
		r := tree.r
		tree.r = nil
		r.l = tree
		tree = r
		if flag {
			return insert(nums1[1:], nums2, tree)
		}
		return insert(nums1, nums2[1:], tree)
	}

	panic(fmt.Sprintf("d(tree.l) < d(tree.r) occurred: %d, %d", d(tree.l), d(tree.r)))
}




