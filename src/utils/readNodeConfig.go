package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type NodeGroupConfig struct {
    Name                 string `yaml:"name" validate:"required"`
    ScalingConfiguration ScalingConfig `yaml:"scalingConfiguration" validate:"required"`
    NetworkConfiguration NetworkConfig `yaml:"networkConfiguration" validate:"required"`
    RoleArn              string `yaml:"roleArn" validate:"omitempty,rolearn"`
    ComputeConfiguration ComputeConfig `yaml:"computeConfiguration" validate:"required"`
    Tags                 map[string]string `yaml:"tags" validate:"required,dive"`
    KubernetesLabels     map[string]string `yaml:"kubernetesLabels" validate:"required,dive"`
    KubernetesTaints     []KubernetesTaint `yaml:"kubernetesTaints" validate:"required,dive"`
}

type ScalingConfig struct {
    DesiredCapacity     int `yaml:"desiredCapacity" validate:"required,minfield=MinSize,maxfield=MaxSize"`
    MinSize             int `yaml:"minSize" validate:"required,min=1"`
    MaxSize             int `yaml:"maxSize" validate:"required,minfield=MinSize"`
    MaximumUnavailable  MaximumUnavailable `yaml:"maximumUnavailable" validate:"required"`
}

type MaximumUnavailable struct {
    Type  string `yaml:"type" validate:"required"`
    Value int    `yaml:"value" validate:"required"`
}

type NetworkConfig struct {
    SubnetIds         []string `yaml:"subnetIds" validate:"required,dive,subnetid"`
    Ec2KeyPair        string `yaml:"ec2KeyPair" validate:"required"`
    SecurityGroupIds  []string `yaml:"securityGroupIds" validate:"required,dive,securitygroupid"`
}

type ComputeConfig struct {
    AmiType        string   `yaml:"amiType" validate:"required"`
    CapacityType   string   `yaml:"capacityType" validate:"required"`
    InstanceTypes  []string `yaml:"instanceTypes" validate:"required,dive,instancetype"`
    DiskSize       int      `yaml:"diskSize" validate:"required"`
}

type KubernetesTaint struct {
    Key    string `yaml:"key" validate:"required"`
    Value  string `yaml:"value" validate:"required"`
    Effect string `yaml:"effect" validate:"required"`
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