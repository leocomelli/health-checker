package core

import (
	"io/ioutil"

	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v2"
)

const ConfigFilename string = "health.yml"

type Checker interface {
	Check(c echo.Context) error
}

type Health struct {
	Services []Service `yaml:"health"`
}

type Service struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func (h Health) GetByType(name string) []Service {
	var r []Service
	for _, s := range h.Services {
		if s.Name == name {
			r = append(r, s)
		}
	}
	return r
}

func LoadServices() (*Health, error) {
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
