package lc

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

func TestAddTwoNumbers1(t *testing.T) {
	o := NewGomegaWithT(t)

	l1 := &ListNode{2, &ListNode{4, &ListNode{3, nil}}}
	l2 := &ListNode{5, &ListNode{6, &ListNode{4, nil}}}

	fmt.Printf("l1: %d\n", ln2number(l1))
	fmt.Printf("l2: %d\n", ln2number(l2))

	r := addTwoNumbers(l1, l2)
	fmt.Printf("r: %d\n", ln2number(r))

	o.Expect(ln2number(r)).To(Equal(ln2number(l1) + ln2number(l2)))
}

func ln2number(l *ListNode) int {
	if l == nil {
		return 0
	}

	return l.Val + 10*ln2number(l.Next)
}


func TestAddTwoNumbers2(t *testing.T) {
	o := NewGomegaWithT(t)

	l1 := &ListNode{5, nil}
	l2 := &ListNode{5, nil}

	fmt.Printf("l1: %d\n", ln2number(l1))
	fmt.Printf("l2: %d\n", ln2number(l2))

	r := addTwoNumbers(l1, l2)
	fmt.Printf("r: %d\n", ln2number(r))

	o.Expect(ln2number(r)).To(Equal(ln2number(l1) + ln2number(l2)))
}

func TestAddTwoNumbers3(t *testing.T) {
	o := NewGomegaWithT(t)

	l1 := &ListNode{4, nil}
	l2 := &ListNode{5, nil}

	fmt.Printf("l1: %d\n", ln2number(l1))
	fmt.Printf("l2: %d\n", ln2number(l2))

	r := addTwoNumbers(l1, l2)
	fmt.Printf("r: %d\n", ln2number(r))

	o.Expect(ln2number(r)).To(Equal(ln2number(l1) + ln2number(l2)))
}