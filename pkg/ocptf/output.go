package ocptf

import (
	"fmt"
)

const (
	Master           = "master"
	UnderlineMetaKey = "_meta"
	HostVarsKey      = "hostvars"
)

type ListOutput struct {
	GroupMap map[string]interface{}
}

type HostOutput struct {
	VarMap map[string]interface{}
}

type GroupI interface {
	GetGroupMap() map[string]interface{}
}

type Group struct {
	Name     string                 `json:"-"`
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars"`
	Children []string               `json:"children"`
}

type Host struct {
	Name   string
	VarMap map[string]interface{}
}

func GetListOutput(groups []Group, hosts []Host) (*ListOutput, error) {
	groupMap := make(map[string]interface{})
	for _, g := range groups {
		if g.Name == UnderlineMetaKey {
			return nil, fmt.Errorf("key (%s) is reserved", UnderlineMetaKey)
		}
		if groupMap[g.Name] != nil {
			return nil, fmt.Errorf("key (%s) exists already in groupMap with value (%v)", g.Name, groupMap[g.Name])
		}
		groupMap[g.Name] = g
	}
	HostVars := make(map[string]map[string]interface{})
	hostVarsMap := make(map[string]interface{})
	for _, h := range hosts {
		if hostVarsMap[h.Name] != nil {
			return nil, fmt.Errorf("key (%s) exists already in hostVarsMap with value (%v)", h.Name, hostVarsMap[h.Name])
		}
		hostVarsMap[h.Name] = h.VarMap
	}
	HostVars[HostVarsKey] = hostVarsMap
	groupMap[UnderlineMetaKey] = HostVars
	listOutput := ListOutput{groupMap}
	return &listOutput, nil

}
