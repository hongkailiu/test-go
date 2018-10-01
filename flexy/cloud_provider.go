package flexy

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

var (
	DryRunnerCounter = 0
)

type CloudProvider interface {
	CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error
	WaitUntilRunning(host *Host, timeout time.Duration) error
}

type AWS struct {
	SVC *ec2.EC2
}

func (aws *AWS) CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error {
	name := configParams["name"]
	imageID := configParams["imageID"]
	kubernetesClusterValue := configParams["kubernetesClusterValue"]

	instances, err := CreateInstancesOnEC2(aws.SVC, name, imageID, int64(1), role.EC2Params.InstanceType, kubernetesClusterValue, role.EC2Params.BlockDeviceMappings)
	if err != nil {
		return err
	}
	if len(instances) != 1 {
		return errors.New(fmt.Sprintf("NOT 1 instance: %d", len(instances)))
	}
	instance := instances[0]
	log.WithFields(log.Fields{"instance.InstanceId": *instance.InstanceId,}).Info("instance created")
	(*host).ID = *instance.InstanceId
	return nil
}

func (aws *AWS) WaitUntilRunning(host *Host, timeout time.Duration) error {
	return WaitUntilRunningOnEC2(aws.SVC, (*host).ID, timeout, host)
}

type DryRunner struct {
}

func (dr *DryRunner) CreateAnInstance(role OCPRole, configParams map[string]string, host *Host) error {
	host.ID = fmt.Sprintf("dry-run-%s", uuid.New().String())
	host.IPv4PublicIP = fmt.Sprintf("23.23.23.%d", DryRunnerCounter+23)
	host.PublicDNS = fmt.Sprintf("dry-run-23-23-23-%d.example.com", DryRunnerCounter+23)
	DryRunnerCounter++
	return nil
}

func (dr *DryRunner) WaitUntilRunning(host *Host, timeout time.Duration) error {
	return nil
}
