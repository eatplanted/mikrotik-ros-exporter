package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Credential struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Configuration struct {
	Timeout     float64
	Credentials map[string]Credential `yaml:"credentials"`
}

func NewConfiguration(configFilePath string) (Configuration, error) {
	content, err := os.ReadFile(configFilePath)
	if err != nil {
		return Configuration{}, err
	}

	var config Configuration
	if err := yaml.Unmarshal(content, &config); err != nil {
		return Configuration{}, err
	}

	return config, nil
}

func (c Configuration) FindCredential(name string) (Credential, error) {
	if credential, ok := c.Credentials[name]; ok {
		return credential, nil
	}

	return Credential{}, errors.New("credential not found")
}
