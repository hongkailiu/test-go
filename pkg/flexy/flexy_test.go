package flexy

import (
	"fmt"

	"github.com/hongkailiu/test-go/flexy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var _ = Describe("[Main] Flexy", func() {

	var (
		inputPath            = "../test_files/flexy/template.yaml"
		inputInventoryFolder = "../test_files/flexy/inv"
		outputFolder         = "../build/output/flexy/inv"
		cp                   flexy.CloudProvider
	)

	Context("With an input folder", func() {
		It("should start flexy job", func() {
			log.WithFields(log.Fields{"inputPath": inputPath}).Debug("loading config from path:")
			config := flexy.OCPClusterConfig{}
			Expect(flexy.LoadOCPClusterConfig(inputPath, &config)).To(BeNil())
			log.WithFields(log.Fields{"config.CloudProvider": config.CloudProvider}).Debug("config.CloudProvider found:")
			switch config.CloudProvider {
			case flexy.CloudProviderAWS:
				cp = flexy.AWS{SVC: svc}
			case flexy.CloudProviderDryRunner:
				cp = flexy.DryRunner{}
			case flexy.CloudProviderGCE:
				cp = flexy.GCE{SVC: computeService}
				Skip("Does not support GCE yet")
			default:
				Fail(fmt.Sprintf("The required cloud provider is not implemented: %s", config.CloudProvider))
			}

			bytes, _ := yaml.Marshal(config)
			fmt.Println("=========\n" + string(bytes))

			Expect(flexy.Start(cp, config, inputInventoryFolder, outputFolder)).To(BeNil())
		})
	})
})
