package utils_test

import (
	"testing"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
    t.Parallel()
    // Test case when a required field is missing
    t.Run("TestClusterConfigNameIsMissing", func(t *testing.T) {
        config := utils.ClusterConfig{
            // Don't set the Name field
            Version:           "1.18",
        }
        err := utils.ValidateConfigs(config)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Name is required")
    })

    // Test case when a required field is empty
    t.Run("TestClusterConfigNameIsEmpty", func(t *testing.T) {
        config := utils.ClusterConfig{
            Name:             "",
            Version:           "1.18",
        }
        err := utils.ValidateConfigs(config)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Name cannot be empty")
    })

    // Test case when a field is out of range
    t.Run("TestNodeGroupMinCountIsGreaterThanMaxCount", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         2,
                MaxSize:         1,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "MaxSize must be greater than MinSize")
    })

    // Test case when DesiredCapacity is out of range
    t.Run("TestNodeGroupDesiredCountIsGreaterThanMaxCount", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 3,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "DesiredCapacity must be between MinSize and MaxSize")
    })

    // Test case when a field is invalid
    t.Run("TestNodeGroupInstanceTypeIsInvalid", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"invalid-instance-type"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "InstanceType must be a valid EC2 instance type")
    })

    // Test case when a field is missing
    t.Run("TestNodeGroupRoleArnIsMissing", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "RoleArn must be specified")
    })

    // Test case when a field is empty
    t.Run("TestNodeGroupRoleArnIsEmpty", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:        "",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "RoleArn cannot be empty")
    })

    // Test case when a field is invalid
    t.Run("TestNodeGroupRoleArnIsInvalid", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:        "invalid-role-arn",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:           map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "RoleArn must be a valid IAM role ARN")
    })

    // Test case when a field is missing
    t.Run("TestNodeGroupTagsAreMissing", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Tags must be specified")
    })

    // Test case when a field is empty
    t.Run("TestNodeGroupTagsAreEmpty", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags: nil,
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Tags cannot be empty")
    })

    t.Run("TestValidSubnetID", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:          map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "SubnetIds must be valid AWS subnet IDs")
    })

    t.Run("TestInvalidSubnetID", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:          map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_SCHEDULE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "SubnetIds must be valid AWS subnet IDs")
    })

    t.Run("TestValidSecurityGroupID", func(t *testing.T) {
        cluster := utils.ClusterConfig{
            Name:             "my-cluster",
            Version:           "1.18",
            RoleArn:          "arn:aws:iam::123456789012:role/eks-cluster-role",
            SubnetIds:        []string{"subnet-12345678912345678"},
            SecurityGroupIds: []string{"sg-0f3a7d6b8e5c4e5c9"},
            PublicAccessCidrs: []string{"10.0.0.0/16"},
            Tags:            map[string]string{"key": "value"},
            ServiceIpv4Cidr: "172.20.0.0/16",
        }
        err := utils.ValidateConfigs(cluster)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "SecurityGroupIds must be valid AWS security group IDs")
    })

    t.Run("TestInvalidSecurityGroupID", func(t *testing.T) {
        cluster := utils.ClusterConfig{
            Name:             "my-cluster",
            Version:           "1.18",
            RoleArn:          "arn:aws:iam::123456789012:role/eks-cluster-role",
            SubnetIds:        []string{"subnet-027691384e95e1c10"},
            SecurityGroupIds: []string{"sg-12345"},
            PublicAccessCidrs: []string{"10.0.0.0/16"},
            Tags:            map[string]string{"key": "value"},
        }
        err := utils.ValidateConfigs(cluster)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "SecurityGroupIds must be valid AWS security group IDs")
    })
            

    t.Run("TestInvalidKubernetesTaintEffect", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:          map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "INVALID_EFFECT",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Effect must be one of NoSchedule, PreferNoSchedule, NoExecute")
    })

    t.Run("TestValidKubernetesTaintEffect", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 1,
                MinSize:         1,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:          map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_EXECUTE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "Effect must be one of NoSchedule, PreferNoSchedule, NoExecute")
    })

    t.Run("TestZeroMinZeroDesiredCapacity", func(t *testing.T) {
        nodeGroup := utils.NodeGroupConfig{
            Name:            "my-node-group",
            ScalingConfiguration: utils.ScalingConfig{
                DesiredCapacity: 0,
                MinSize:         0,
                MaxSize:         2,
                MaximumUnavailable: utils.MaximumUnavailable{
                    Type:  "percentage",
                    Value: 50,
                },
            },
            NetworkConfiguration: utils.NetworkConfig{
                SubnetIds:      []string{"subnet-12345678912345678"},
                Ec2KeyPair:     "my-key-pair",
                SecurityGroupIds: []string{"sg-12345678912345678"},
            },
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            ComputeConfiguration: utils.ComputeConfig{
                AmiType:        "AL2_x86_64",
                CapacityType:   "ON_DEMAND",
                InstanceTypes:  []string{"t3.medium"},
                DiskSize:       20,
            },
            Tags:          map[string]string{"key": "value"},
            KubernetesLabels: map[string]string{"key": "value"},
            KubernetesTaints: []utils.KubernetesTaint{
                {
                    Key:    "key",
                    Value:  "value",
                    Effect: "NO_EXECUTE",
                },
            },
        }
        err := utils.ValidateConfigs(nodeGroup)
        require.NoError(t, err)
        // require.Contains(t, err.Error(), "DesiredCapacity is between MinSize and MaxSize")
    })
}
