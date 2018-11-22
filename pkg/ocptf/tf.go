package ocptf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// VERSION of the ocptf cmd
	VERSION = "0.0.1"
)

func LoadTFStateFile(path string) (*TFState, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tfState := TFState{}
	err = json.Unmarshal(bytes, &tfState)
	if err != nil {
		return nil, err
	}
	return &tfState, nil
}

type TFState struct {
	Modules []Module
}

type Module struct {
	Resources map[string]Resource
}

type Resource struct {
	Type    string
	Primary Primary
}

type Primary struct {
	ID         string
	Attributes map[string]string
}

func load(path string) ([]Group, []Host, error) {
	tfState, err := LoadTFStateFile(path)
	if err != nil {
		return nil, nil, err
	}
	if len(tfState.Modules) != 1 {
		return nil, nil, fmt.Errorf("len(tfState.Modules) is %d", len(tfState.Modules))
	}

	osev3Group := Group{Name: "OSEv3", Vars: map[string]interface{}{}, Children: []string{"masters", "nodes", "etcd"}, Hosts: []string{}}
	mastersGroup := Group{Name: "masters", Vars: map[string]interface{}{}, Children: []string{}}
	nodesGroup := Group{Name: "nodes", Vars: map[string]interface{}{}, Children: []string{}}
	etcdGroup := Group{Name: "etcd", Vars: map[string]interface{}{}, Children: []string{}}
	glusterGroup := Group{Name: "glusterfs", Vars: map[string]interface{}{}, Children: []string{}}
	var hosts []Host
	for k, v := range tfState.Modules[0].Resources {
		log.WithFields(log.Fields{"k": k, "v": v}).Debug("resource")

		if v.Type != "aws_instance" {
			continue
		}
		if osev3Group.Vars["openshift_cloudprovider_kind"] == nil {
			osev3Group.Vars["openshift_cloudprovider_kind"] = "aws"
		}
		h := Host{Name: v.Primary.Attributes["public_dns"], ID: v.Primary.ID, PublicIP: v.Primary.Attributes["public_ip"]}
		h.VarMap = map[string]interface{}{"openshift_public_hostname": h.Name}

		if !strings.Contains(k, "master") && !strings.Contains(k, "infra") &&
			!strings.Contains(k, "node") && !strings.Contains(k, "worker") &&
			!strings.Contains(k, "etcd") && !strings.Contains(k, "gluster") {
			return nil, nil, fmt.Errorf("malformed instance name,: %s", k)
		}

		if strings.Contains(k, "master") {
			mastersGroup.Hosts = append(mastersGroup.Hosts, h.Name)
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = "node-config-master"
			h.VarMap["openshift_schedulable"] = true
		}
		if strings.Contains(k, "infra") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = "node-config-infra"
			if osev3Group.Vars["openshift_master_default_subdomain"] == nil {
				osev3Group.Vars["openshift_master_default_subdomain"] = fmt.Sprintf("apps.%s.xip.io", h.PublicIP)
			}
		}
		if strings.Contains(k, "worker") || strings.Contains(k, "node") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = "node-config-compute"
		}
		if strings.Contains(k, "gluster") {
			nodesGroup.Hosts = append(nodesGroup.Hosts, h.Name)
			glusterGroup.Hosts = append(glusterGroup.Hosts, h.Name)
			h.VarMap["openshift_node_group_name"] = "node-config-compute"
			h.VarMap["glusterfs_devices"] = `'["/dev/nvme2n1"]'`
		}

		hosts = append(hosts, h)
	}

	if len(etcdGroup.Hosts) == 0 {
		etcdGroup.Hosts = mastersGroup.Hosts
	}

	log.WithFields(log.Fields{"hosts": hosts}).Debug("hosts")

	groups := []Group{osev3Group, mastersGroup, nodesGroup, etcdGroup, glusterGroup}
	if strings.ToLower(os.Getenv("install_ocp_gluster")) == "false" {
		groups = []Group{osev3Group, mastersGroup, nodesGroup, etcdGroup}
	}

	return groups, hosts, nil
}

func DoList(path string, dynamic bool) error {
	log.WithFields(log.Fields{"path": path, "dynamic": dynamic}).Debug("DoList")
	groups, hosts, err := load(path)
	if err != nil {
		return err
	}
	listOutput, err := GetListOutput(groups, hosts)
	if err != nil {
		return err
	}

	if dynamic {
		bytes, err := json.MarshalIndent(listOutput.GroupMap, "", "  ")
		if err != nil {
			return err
		}
		jsonString := fmt.Sprintf(string(bytes))
		log.WithFields(log.Fields{"jsonString": jsonString}).Debug("")
		fmt.Println(jsonString)

	} else {
		return fmt.Errorf("TODO: static inv output")
	}

	return nil
}

func DoHost(path, name string, dynamic bool) error {
	_, hosts, err := load(path)
	if err != nil {
		return err
	}
	hostOutput := GetHostOutput(name, hosts)

	if dynamic {
		bytes, err := json.MarshalIndent(hostOutput.VarMap, "", "  ")
		if err != nil {
			return err
		}
		jsonString := fmt.Sprintf(string(bytes))
		fmt.Println(jsonString)

	} else {
		return fmt.Errorf("TODO: static inv output")
	}

	return nil
}
