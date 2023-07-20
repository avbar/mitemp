package config

import (
	"os"

	"github.com/avbar/mitemp/internal/handler"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Sensors []handler.Sensor `yaml:"sensors"`
	Port    int              `yaml:"port"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing config file")
	}

	return nil
}
