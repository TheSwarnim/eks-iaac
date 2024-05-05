package utils

import (
	"log"
	"os"
	"path/filepath"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type ClusterConfig struct {
    Name              string `yaml:"name" validate:"required"`
    Region            string `yaml:"region" validate:"required"`
    Version           string `yaml:"version" validate:"required"`
    RoleArn           string `yaml:"roleArn" validate:"required"`
    PublicAccessCidrs []string `yaml:"publicAccessCidrs" validate:"required"`
    SecurityGroupIds  []string `yaml:"securityGroupIds" validate:"required"`
    SubnetIds         []string `yaml:"subnetIds" validate:"required"`
    Tags              map[string]string `yaml:"tags" validate:"required"`
    NodeGroups        []NodeGroup `yaml:"nodeGroups" validate:"required"`
}

type NodeGroup struct {
	Name            string `yaml:"name" validate:"required"`
    InstanceType    string `yaml:"instanceType" validate:"required"`
    DesiredCapacity int    `yaml:"desiredCapacity" validate:"required"`
    MinSize         int    `yaml:"minSize" validate:"required"`
    MaxSize         int    `yaml:"maxSize" validate:"required"`
	SubnetIds       []string `yaml:"subnetIds"`
	RoleArn         string `yaml:"roleArn" validate:"required"`
	Tags 		  map[string]string `yaml:"tags" validate:"required"`
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

		validate := validator.New()
		err = validate.Struct(clusterConfig)
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
