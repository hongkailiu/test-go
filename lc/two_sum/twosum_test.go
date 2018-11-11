package two_sum

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test1(t *testing.T) {
	o := NewGomegaWithT(t)

	r := twoSum(nil, 0)

	o.Expect(r).To(BeNil())
}
