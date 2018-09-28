package flexy

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

func Start(svc *ec2.EC2, config OCPClusterConfig, inputPath string, outputFolder string) error {
	if len(config.OCPRoles) > 6 {
		return errors.New(fmt.Sprintf("too many roles: %d", len(config.OCPRoles)))
	}

	var masterGroup, etcdGroup, nodeGroup, lbGroup, glusterFSGroup []Host
	ocVars := map[string]string{}

	for _, role := range config.OCPRoles {
		if err := ValidateOCPRole(role.Name); err != nil {
			return err
		}
		if role.Size > 25 {
			return errors.New(fmt.Sprintf("too many instance for role %s: %d", role.Name, role.Size))
		}

		if role.Size > 0 {
			for i := 1; i <= role.Size; i++ {
				instances, err := CreateInstances(svc, fmt.Sprintf("%s-%s-%d", config.InstancePrefix, role.Name, i),
					config.ImageID, int64(1), role.InstanceType, config.KubernetesClusterValue, role.BlockDeviceMappings)
				if err != nil {
					return err
				}

				if len(instances) != 1 {
					return errors.New(fmt.Sprintf("NOT 1 instance: %d", len(instances)))
				}

				instance := instances[0]
				log.WithFields(log.Fields{"instance.InstanceId": *instance.InstanceId,}).Info("instance created")
				host := Host{}
				host.ID = *instance.InstanceId
				host.OCMasterSchedulable = false
				switch role.Name {
				case OCPRoleMaster:
					host.OCNodeGroupName = "node-config-master"
					host.OCMasterSchedulable = true
				case OCPRoleInfra:
					host.OCNodeGroupName = "node-config-infra"
				case OCPRoleCompute:
					host.OCNodeGroupName = "node-config-compute"
				}

				err = WaitUntilRunning(svc, *instance.InstanceId, 2*time.Minute, &host)
				if err != nil {
					return err
				}
				if role.Name == OCPRoleInfra && len(ocVars["openshift_master_default_subdomain"]) == 0 {
					//"openshift_master_default_subdomain": "apps.someip.xip.io",
					ocVars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", host.IPv4PublicIP)
				}

				switch role.Name {
				case OCPRoleMaster:
					masterGroup = append(masterGroup, host)
					nodeGroup = append(nodeGroup, host)
				case OCPRoleInfra:
					nodeGroup = append(nodeGroup, host)
				case OCPRoleCompute:
					nodeGroup = append(nodeGroup, host)
				}


			}

		}

	}
	if len(etcdGroup) == 0 {
		etcdGroup = masterGroup
	}
	inventoryConfig := Config{masterGroup, etcdGroup, nodeGroup, lbGroup, glusterFSGroup, ocVars}
	err := Generate(inputPath, inventoryConfig, outputFolder)
	if err != nil {
		return err
	}
	return nil
}
