package flexy_test

import (
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	svc *ec2.EC2
)

func TestFlexy(t *testing.T) {

	RegisterFailHandler(Fail)
	log.SetLevel(log.DebugLevel)
	BeforeSuite(func() {
		log.Debug("BeforeSuite============")
		log.Debug("here i am")
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			panic("unable to load SDK config, " + err.Error())
		}

		// Set the AWS Region that the service clients should use
		cfg.Region = endpoints.UsWest2RegionID

		// Using the Config value, create the DynamoDB client
		svc = ec2.New(cfg)

	})
	RunSpecs(t, "Flexy Suite")
}
