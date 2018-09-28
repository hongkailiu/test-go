package flexy_test

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hongkailiu/test-go/flexy"
	log "github.com/sirupsen/logrus"
)



var _ = Describe("AwsEc2Helper", func() {
	var (
		imageID                string
		count                  int64
		instanceType           ec2.InstanceType
		kubernetesClusterValue string
		blockDeviceMappings    []ec2.BlockDeviceMapping
		svc                    *ec2.EC2
	)



	BeforeEach(func() {
		log.Debug("BeforeEach============")
		imageID = "ami-02a29c6a327959825"
		count = int64(1)
		instanceType = ec2.InstanceTypeM5Large
		kubernetesClusterValue = "hongkliu-kc"
		blockDeviceMappings = []ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeType: ec2.VolumeTypeGp2,
					VolumeSize: aws.Int64(66),
				},
			},
		}
	})

	Context("With an ec2 instance", func() {
		It("should create instance", func() {
			log.Info(imageID, count, instanceType, kubernetesClusterValue, blockDeviceMappings)
			//instances, err := flexy.CreateInstances(svc, "hongkliu-311-001", imageID, count, instanceType, kubernetesClusterValue, blockDeviceMappings)
			//Expect(err).To(BeNil())
			//Expect(instances).NotTo(BeNil())
			//for _, instance := range instances {
			//	log.WithFields(log.Fields{"instance.InstanceId": *instance.InstanceId,}).Debug("instance created")
			//}
			//id := *(instances[0].InstanceId)
			id := "i-012002d42f46bf4f0"
			host := flexy.Host{}
			Expect(flexy.WaitUntilRunning(svc, id, 2*time.Minute, &host)).To(BeNil())

		})
	})

})
