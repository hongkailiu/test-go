package flexy_test

import (
	"fmt"
	"github.com/hongkailiu/test-go/flexy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

var _ = Describe("[Main] Flexy", func() {

	var (
		inputPath    = "../test_files/flexy/template.yaml"
		inputInventoryFolder = "../test_files/flexy/inv"
		outputFolder = "../build/output/flexy/inv"
	)

	Context("With an input folder", func() {
		It("should start flexy job", func() {
			config := flexy.OCPClusterConfig{}
			Expect(flexy.LoadOCPClusterConfig(inputPath, &config)).To(BeNil())
			Expect(config.OCPRoles).Should(HaveLen(3))

			bytes, _ := yaml.Marshal(config)
			fmt.Println("=========\n" + string(bytes))

			Expect(flexy.Start(svc, config, inputInventoryFolder, outputFolder)).To(BeNil())
		})
	})
})
