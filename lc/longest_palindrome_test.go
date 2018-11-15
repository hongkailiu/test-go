package lc

import (
	"testing"

	. "github.com/onsi/gomega"
)

// https://stackoverflow.com/questions/16935965/how-to-run-test-cases-in-a-specified-file
// go test -run TestLongestPalindrome1 ./lc/...
func TestLongestPalindrome1(t *testing.T) {
	o := NewGomegaWithT(t)

	r := longestPalindrome("babad")

	o.Expect(r).Should(Or(Equal("bab"), Equal("aba")))
}

func TestLongestPalindrome2(t *testing.T) {
	o := NewGomegaWithT(t)

	r := longestPalindrome("cbbd")

	o.Expect(r).Should(Equal("bb"))
}

func TestLongestPalindrome3(t *testing.T) {
	o := NewGomegaWithT(t)

	r := longestPalindrome("a")

	o.Expect(r).Should(Equal("a"))
}


func TestLongestPalindrome4(t *testing.T) {
	o := NewGomegaWithT(t)

	r := longestPalindrome("bb")

	o.Expect(r).Should(Equal("bb"))
}

func TestLongestPalindrome5(t *testing.T) {
	o := NewGomegaWithT(t)

	r := longestPalindrome("abcda")

	o.Expect(r).Should(Equal("a"))
}