package config

import (
	"log"
	"os"

	"api-registry/pkg/constants"
	"api-registry/pkg/model"

	"gopkg.in/yaml.v2"
)

func InitConfig() {

	configPath := constants.DefaultConfigFile

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var config model.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	model.SetConfig(config)

}
