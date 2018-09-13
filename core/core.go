package core

import (
	"io/ioutil"
	"os"

	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v2"
)

var ConfigFilename string = "health.yml"

type Checker interface {
	Check(c echo.Context) error
}

type Health struct {
	Services []Service `yaml:"health"`
}

type Service struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func (h Health) GetByType(tp string) []Service {
	var r []Service
	for _, s := range h.Services {
		if s.Type == tp {
			r = append(r, s)
		}
	}
	return r
}

func LoadServices() (*Health, error) {
	customFilepath := os.Getenv("HC_FILE")
	if customFilepath != "" {
		ConfigFilename = customFilepath
	}

	data, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		return nil, err
	}

	var t Health
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
