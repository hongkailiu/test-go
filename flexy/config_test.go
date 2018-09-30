package flexy_test

import (
	"github.com/hongkailiu/test-go/flexy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		inputPath = "../test_files/flexy/unit.test.files/template.yaml"
	)

	Context("With an input folder", func() {
		It("should load the config file", func() {
			config := flexy.OCPClusterConfig{}
			Expect( flexy.LoadOCPClusterConfig(inputPath, &config)).To(BeNil())
			Expect(config.OCPRoles).Should(HaveLen(3))
		})
	})

})
