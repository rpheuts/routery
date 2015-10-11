package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"log"
)

func GetConfig(configPath string) *RouteryConfig {
	out, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
		return nil
	}

	config := RouteryConfig{}
	err = yaml.Unmarshal(out, &config)

	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
		return nil
	}

	return &config
}