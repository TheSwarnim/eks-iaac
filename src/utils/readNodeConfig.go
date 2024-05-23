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
    KubernetesTaints     []KubernetesTaint `yaml:"kubernetesTaints" validate:"omitempty,dive"`
}

type ScalingConfig struct {
    DesiredCapacity     int `yaml:"desiredCapacity" validate:"minfield=MinSize,maxfield=MaxSize"`
    MinSize             int `yaml:"minSize" validate:"maxfield=DesiredCapacity,min=0"`
    MaxSize             int `yaml:"maxSize" validate:"minfield=DesiredCapacity,min=1"`
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
    DiskSize       int      `yaml:"diskSize" validate:"required,min=8"`
}

type KubernetesTaint struct {
    Key    string `yaml:"key" validate:"required"`
    Value  string `yaml:"value" validate:"required"`
    Effect string `yaml:"effect" validate:"required,tainteffect"`
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