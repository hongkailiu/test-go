package weather

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppID  string   `json:"appID"`
	Cities []City   `json:"cities"`
	Writers []string `json:"writer"`
}

type City struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func LoadConfig(path string) (Config, error) {
	var c Config
	content, err := ioutil.ReadFile(path)
	logrus.WithField("string(content)", string(content)).Debug("file content")
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
