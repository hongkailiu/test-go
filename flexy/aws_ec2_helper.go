package flexy

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"k8s.io/apimachinery/pkg/util/wait"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	nameKey              = "Name"
	securityGroupId      = "sg-5c5ace38"
	keyNameValue         = "id_rsa_perf"
	subnetId             = "subnet-4879292d"
	kubernetesClusterKey = "KubernetesCluster"
)

func CreateInstancesOnEC2(svc *ec2.EC2, name string, imageID string, count int64, instanceType ec2.InstanceType, kubernetesClusterValue string, blockDeviceMappings []ec2.BlockDeviceMapping) ([]ec2.Instance, error) {
	log.WithFields(log.Fields{"name": name,}).Info("instance creating")
	req := svc.RunInstancesRequest(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String(imageID),
		InstanceType: instanceType,
		MinCount:     aws.Int64(count),
		MaxCount:     aws.Int64(count),
		SecurityGroupIds: []string{
			securityGroupId,
		},
		KeyName:             aws.String(keyNameValue),
		SubnetId:            aws.String(subnetId),
		BlockDeviceMappings: blockDeviceMappings,
		TagSpecifications: []ec2.TagSpecification{

			{
				ResourceType: ec2.ResourceTypeInstance,
				Tags: []ec2.Tag{
					{
						Key:   &nameKey,
						Value: &name,
					},
					{
						Key:   &kubernetesClusterKey,
						Value: &kubernetesClusterValue,
					},
				},
			},
		},
	})

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp.Instances, nil
}

func WaitUntilRunningOnEC2(svc *ec2.EC2, instanceId string, timeout time.Duration, host *Host) error {
	err := wait.Poll(10*time.Second, timeout,
		func() (bool, error) {
			log.Debug(fmt.Sprintf("checking if the instance %s is running ...", instanceId))
			instance, err := DescribeAInstance(svc, instanceId)
			if err != nil {
				return false, nil
			}
			log.WithFields(log.Fields{"instance.State.Code": strconv.FormatInt(*instance.State.Code, 10),}).Debug("instance code found")
			if *(instance.State.Code) == int64(16) {
				host.PublicDNS = *instance.PublicDnsName
				host.IPv4PublicIP = *instance.PublicIpAddress
				return true, nil
			}
			return false, nil
		})
	return err
}

func DescribeAInstance(svc *ec2.EC2, instanceId string) (*ec2.Instance, error) {
	req := svc.DescribeInstancesRequest(&ec2.DescribeInstancesInput{
		InstanceIds: []string{
			instanceId,
		},
	})

	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	if resp.Reservations == nil || resp.Reservations[0].Instances == nil {
		return nil, errors.New(fmt.Sprintf("no found instance by id: %s", instanceId))
	}

	return &resp.Reservations[0].Instances[0], nil
}
