package ocptf_test

import (
	"testing"

	. "github.com/hongkailiu/test-go/pkg/ocptf"
	. "github.com/onsi/gomega"
)

func TestLoadTFStateFile1(t *testing.T) {
	o := NewGomegaWithT(t)

	r, err := LoadTFStateFile("../../test_files/ocpft/unit.test.files/terraform.tfstate.json")
	o.Expect(err).To(BeNil())
	o.Expect(r).NotTo(BeNil())

}

func TestStart1(t *testing.T) {
	o := NewGomegaWithT(t)

	err := DoList("../../test_files/ocpft/unit.test.files/terraform.tfstate.json", true)
	o.Expect(err).To(BeNil())

}

func TestStart2(t *testing.T) {
	o := NewGomegaWithT(t)

	err := DoList("../../test_files/ocpft/unit.test.files/terraform.tfstate.gfs.json", true)
	o.Expect(err).To(BeNil())

}

func TestStart3(t *testing.T) {
	o := NewGomegaWithT(t)

	err := DoList("../../test_files/ocpft/unit.test.files/terraform.tfstate.json", false)
	o.Expect(err).To(BeNil())

}

func TestStart4(t *testing.T) {
	o := NewGomegaWithT(t)

	err := DoList("../../test_files/ocpft/unit.test.files/terraform.tfstate.gfs.json", false)
	o.Expect(err).To(BeNil())

}
