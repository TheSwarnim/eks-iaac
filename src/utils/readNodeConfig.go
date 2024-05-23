package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type NodeGroupConfig struct {
	Name            string `yaml:"name" validate:"required"`
    InstanceType    string `yaml:"instanceType" validate:"required,instancetype"`
    DesiredCapacity int    `yaml:"desiredCapacity" validate:"required,minfield=MinSize"`
    MinSize         int    `yaml:"minSize" validate:"required,min=1"`
	MaxSize         int    `yaml:"maxSize" validate:"required,minfield=MinSize"`
	SubnetIds       []string `yaml:"subnetIds" validate:"omitempty,dive,subnetid"`
	RoleArn         string `yaml:"roleArn" validate:"omitempty,rolearn"`
	Tags 		  	map[string]string `yaml:"tags" validate:"required"`
}

func ReadNodeConfigs(nodeDirInClusterDir string) ([]NodeGroupConfig, error) {
    var nodeGroupConfigs []NodeGroupConfig

    // Walk through the cluster directory and its subdirectories
    err := filepath.Walk(nodeDirInClusterDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // If the current file is a nodegroup yaml file, read it
        if info.IsDir() || filepath.Ext(path) != ".yaml" {
            return nil
        }

        log.Printf("Reading nodegroup file: %s", path)

        // Read the nodegroup yaml file
        data, err := os.ReadFile(path)
        if err != nil {
            return err
        }

        // Unmarshal the YAML data into a NodeGroup struct
        var nodeGroup NodeGroupConfig
        err = yaml.Unmarshal(data, &nodeGroup)
        if err != nil {
            return err
        }

        // Validate the NodeGroup
        err = ValidateConfigs(&nodeGroup)
        if err != nil {
            return err
        }

        // Add the NodeGroup to the slice
        nodeGroupConfigs = append(nodeGroupConfigs, nodeGroup)

        log.Printf("Successfully read nodegroup file: %s", path)

        return nil
    })

    if err != nil {
        return nil, err
    }

    return nodeGroupConfigs, nil
}