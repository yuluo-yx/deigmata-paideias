package pkg

import (
	"log"
	"os"
	"ospp_search/pkg/model"

	"sigs.k8s.io/yaml"
)

func InitConfig() {

	configPath := "config/config.yml"

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}

	var pro model.Project
	err = yaml.Unmarshal(yamlFile, &pro)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	//log.Printf("config info: %v\n", pro)

	model.SetProject(pro)
}
