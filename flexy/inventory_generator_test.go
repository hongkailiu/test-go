package flexy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hongkailiu/test-go/flexy"
)

var _ = Describe("InventoryGenerator", func() {
	var (
		inputPath    string
		config       flexy.Config
		outputFolder string
	)

	BeforeEach(func() {
		inputPath = "../test_files/flexy/inv"
		config = flexy.Config{
			MasterGroup: []flexy.Host{

						{
							ID:        "111",
							PublicDNS: "001.hongkailiu.tk",
						},
					},
			ETCDGroup:	[]flexy.Host{
						{
							ID:        "111",
							PublicDNS: "001.hongkailiu.tk",
						},
					},
			NodeGroup: []flexy.Host{
						{
							ID:              "001",
							PublicDNS:       "001.hongkailiu.tk",
							OCNodeGroupName: "node-config-master",
							OCMasterSchedulable:   true,
						},
						{
							ID:              "002",
							PublicDNS:       "002.hongkailiu.tk",
							OCNodeGroupName: "node-config-infra",
						},
						{
							ID:              "003",
							PublicDNS:       "003.hongkailiu.tk",
							OCNodeGroupName: "node-config-compute",
						},
						{
							ID:              "004",
							PublicDNS:       "004.hongkailiu.tk",
							OCNodeGroupName: "node-config-compute",
						},
					},
			OCVars: map[string]string{
				"openshift_master_default_subdomain": "apps.someip.xip.io",
			},
		}
		outputFolder = "../build/output/flexy/tmp"

	})

	Context("With an input folder", func() {
		It("should generate the inventory file", func() {
			Expect(flexy.Generate(inputPath, config, outputFolder)).To(BeNil())
		})
	})

	It("should work", func() {

	})
})
