package flexy

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/google/uuid"
)

var (
	COUNTER = 0
)

func CreateInstanceDryrun(name string) ([]ec2.Instance, error) {
	COUNTER = COUNTER + 1
	return []ec2.Instance{
		{
			InstanceId:      aws.String(fmt.Sprintf("dry-run-%s", uuid.New().String())),
			PublicIpAddress: aws.String(fmt.Sprintf("23.23.23.%d", COUNTER+22)),
			PublicDnsName:   aws.String(fmt.Sprintf("dry-run-23-23-23-%d.example.com", COUNTER+22)),
		},
	}, nil
}
