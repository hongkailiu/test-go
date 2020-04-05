package assigner

import (
	"github.com/ghodss/yaml"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	configURL = "https://raw.githubusercontent.com/hongkailiu/test-go/master/cmd/temp/assigner/config.yaml"
)

var (
	client = resty.New()
)

type Service interface {
	GetStatus() (Status, error)
	GetConfig() (Config, error)
	GetGroup(groupName string) ([]string, error)
	SetGroup(groupName string, members []string) error
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}

func (s serviceImpl) GetConfig() (Config, error) {
	c := Config{}
	resp, err := client.R().
		EnableTrace().
		Get(configURL)
	if err != nil {
		return c, err
	}
	logrus.WithField("resp.Body()", string(resp.Body())).Infof("get")
	if err := yaml.Unmarshal(resp.Body(), &c); err != nil {
		return c, err
	}
	return c, nil
}

func (s serviceImpl) SetGroup(groupName string, members []string) error {
	return nil
}

func (s serviceImpl) GetGroup(groupName string) ([]string, error) {
	return []string{"aaa", "bbb"}, nil
}

func (s serviceImpl) GetStatus() (Status, error) {
	status := Status{}
	c, err := s.GetConfig()
	if err != nil {
		return status, err
	}
	status.Config = c
	g, err := s.GetGroup(c.GroupName)
	if err != nil {
		return status, err
	}
	status.Current = g
	return status, nil
}
