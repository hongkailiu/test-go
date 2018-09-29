package flexy

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	log "github.com/sirupsen/logrus"
)

const (
	InstanceCountLimit = 23
	RoleCountLimit     = 6
)

func Start(svc *ec2.EC2, config OCPClusterConfig, inputPath string, outputFolder string) error {
	if len(config.OCPRoles) > RoleCountLimit {
		return errors.New(fmt.Sprintf("RoleCountLimit is %d: too many roles: %d", RoleCountLimit, len(config.OCPRoles)))
	}

	var masterGroup, etcdGroup, nodeGroup, lbGroup, glusterFSGroup []Host
	ocVars := map[string]string{}
	for k, v := range config.OpenshiftAnasibleVar {
		ocVars[k] = v
		ocVars["openshift_clusterid"] = config.KubernetesClusterValue
	}

	instanceCount := 0
	for _, role := range config.OCPRoles {
		if err := ValidateOCPRole(role.Name); err != nil {
			return err
		}
		if instanceCount > InstanceCountLimit {
			return errors.New(fmt.Sprintf("instances over %d", InstanceCountLimit))
		}

		if role.Size > 0 {
			for i := 1; i <= role.Size; i++ {
				name := fmt.Sprintf("%s-%s-%d", config.InstancePrefix, role.Name, i)
				if config.AllInOne {
					name = fmt.Sprintf("%s-%s", config.InstancePrefix, "all-in-one")
				}
				if config.DryRun {
					if instanceCount == 1 {
						return errors.New("required more than 1 instance for all-in-one cluster")
					}
				}
				instances, err := CreateInstanceDryrun(name)
				if !config.DryRun {
					instances, err = CreateInstances(svc, name,
						config.ImageID, int64(1), role.InstanceType, config.KubernetesClusterValue, role.BlockDeviceMappings)
				}
				instanceCount++
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
				case OCPRoleGlusterFS:
					host.OCNodeGroupName = "node-config-compute"
				}

				if !config.DryRun {
					err = WaitUntilRunning(svc, *instance.InstanceId, 2*time.Minute, &host)
					if err != nil {
						return err
					}
				} else {
					host.IPv4PublicIP = *instance.PublicIpAddress
					host.PublicDNS = *instance.PublicDnsName
				}
				if config.AllInOne {
					host.OCNodeGroupName = "node-config-all-in-one"
					ocVars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", host.IPv4PublicIP)
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
				case OCPRoleGlusterFS:
					nodeGroup = append(nodeGroup, host)
					glusterFSGroup = append(glusterFSGroup, host)
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
