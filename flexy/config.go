package flexy

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	MasterGroup    []Host
	ETCDGroup      []Host
	NodeGroup      []Host
	LBGroup        []Host
	GlusterFSGroup []Host
	OCVars         map[string]string
}

const (
	OCPRoleMaster    = "master"
	OCPRoleInfra     = "infra"
	OCPRoleCompute   = "compute"
	OCPRoleGlusterFS = "glusterfs"
)

var (
	OCPRoles = []string{OCPRoleMaster, OCPRoleInfra, OCPRoleCompute, OCPRoleGlusterFS}
)

func ValidateOCPRole(role string) error {
	for _, r := range OCPRoles {
		if role == r {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("invalid OCP role: %s", role))
}

type Host struct {
	ID                  string
	PublicDNS           string
	OCNodeGroupName     string
	OCMasterSchedulable bool
	IPv4PublicIP        string
}
type OCPClusterConfig struct {
	AllInOne               bool              `yaml:"allInOne"`
	DryRun                 bool              `yaml:"dryRun"`
	OCPRoles               []OCPRole         `yaml:"ocpRoles"`
	KubernetesClusterValue string            `yaml:"kubernetesClusterValue"`
	ImageID                string            `yaml:"imageID"`
	InstancePrefix         string            `yaml:"instancePrefix"`
	OpenshiftAnasibleVar   map[string]string `yaml:"openshiftAnasibleVar"`
}

type OCPRole struct {
	Name                string
	Size                int
	InstanceType        ec2.InstanceType         `yaml:"instanceType"`
	CCC                 string                   `yaml:"ccc,omitempty"`
	BlockDeviceMappings []ec2.BlockDeviceMapping `yaml:"blockDeviceMappings"`
}

func LoadOCPClusterConfig(file string, config *OCPClusterConfig) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}

	return nil
}
