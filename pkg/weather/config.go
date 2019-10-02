package weather

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppID     string   `json:"appID"`
	Key       string   `json:"key"`
	Cities    []City   `json:"cities"`
	Writers   []string `json:"writer"`
	OutputDir string   `json:"outputDir"`
}

type City struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	TimeZone string `json:"timezone"`
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

func (c Config) Save(path string) error {
	logrus.WithField("c", c).Debug("saving config")
	bytes, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

func (c Config) Validate() error {
	fileInfo, err := os.Stat(c.OutputDir)
	if os.IsNotExist(err) {
		logrus.WithField("c.OutputDir", c.OutputDir).Info("output dir does not exists, creating ...")
		err := os.MkdirAll(c.OutputDir, 0755)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}
	logrus.WithField("fileInfo.Name()", fileInfo.Name()).Info("Output dir")
	if !fileInfo.IsDir() {
		return fmt.Errorf("output dir '%s' is not a dir", c.OutputDir)
	}
	return nil
}
