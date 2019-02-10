package ocptf_test

import (
	"os"
	"testing"

	. "github.com/hongkailiu/test-go/pkg/ocptf"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

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

func TestStart5(t *testing.T) {
	o := NewGomegaWithT(t)

	err := DoList("../../test_files/ocpft/unit.test.files/terraform.tfstate.all.in.one.json", false)
	o.Expect(err).To(BeNil())

}
