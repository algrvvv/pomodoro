package config

import (
	"os"
	"time"

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

	for idx, i := range Config.Intergations {
		Config.Intergations[idx].Timeout = time.Duration(i.TimeoutInSecond) * time.Second
	}

	return nil
}
