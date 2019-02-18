package flexy_test

import (
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

var (
	svc            *ec2.EC2
	computeService *compute.Service
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
		Expect(svc).NotTo(BeNil())
		data, err := ioutil.ReadFile("/tmp/gce.json")
		if err != nil {
			Fail(err.Error())
		}
		conf, err := google.JWTConfigFromJSON(data, compute.ComputeScope)
		if err != nil {
			Fail(err.Error())
		}
		client := conf.Client(oauth2.NoContext)

		computeService, err = compute.New(client)
		if err != nil {
			Fail(err.Error())
		}

	})
	RunSpecs(t, "Flexy Suite")
}
