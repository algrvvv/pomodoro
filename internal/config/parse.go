package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Parse(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		return err
	}

	return nil
}
