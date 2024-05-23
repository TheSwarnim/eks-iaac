package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ClusterConfig struct {
    Name              string `yaml:"name" validate:"required"`
    Version           string `yaml:"version" validate:"required"`
    RoleArn           string `yaml:"roleArn" validate:"omitempty,rolearn"` 	// roleArn field is optional 
	ServiceIpv4Cidr   string `yaml:"serviceIpv4Cidr" validate:"required,cidrv4"`
    PublicAccessCidrs []string `yaml:"publicAccessCidrs" validate:"required,dive,cidrv4"`
    SecurityGroupIds  []string `yaml:"securityGroupIds" validate:"required,dive,securitygroupid"`
    SubnetIds         []string `yaml:"subnetIds" validate:"required,dive,subnetid"`
    Tags              map[string]string `yaml:"tags" validate:"required,dive"`
}

func ReadClusterConfigs(rootDir string) ([]ClusterConfig, error) {
	var clusterConfigs []ClusterConfig

	// Walk through the root directory and its subdirectories
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the current file is a config.yaml file, read it
		if info.IsDir() || filepath.Base(path) != "config.yaml" {
			return nil
		}

		log.Printf("Reading config file: %s", path)

		// Read the config.yaml file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Unmarshal the YAML data into a ClusterConfig struct
		var clusterConfig ClusterConfig
		err = yaml.Unmarshal(data, &clusterConfig)
		if err != nil {
			return err
		}

		// Validate the ClusterConfig
		err = ValidateConfigs(&clusterConfig)
		if err != nil {
			return err
		}

		// Add the ClusterConfig to the slice
		clusterConfigs = append(clusterConfigs, clusterConfig)

		log.Printf("Successfully read config file: %s", path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return clusterConfigs, nil
}
