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
            // Set other fields as needed...
        }
        err := utils.ValidateClusterConfig(config)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Name is required")
    })

    // Test case when a required field is empty
    t.Run("TestClusterConfigNameIsEmpty", func(t *testing.T) {
        config := utils.ClusterConfig{
            Name:             "",
            Version:           "1.18",
            // Set other fields as needed...
        }
        err := utils.ValidateClusterConfig(config)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Name cannot be empty")
    })

    // Test case when a field is out of range
    t.Run("TestNodeGroupMinCountIsGreaterThanMaxCount", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         3,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            Tags:           map[string]string{"key": "value"},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "MaxSize must be greater than MinSize")
    })

    // Test case when a field is invalid
    t.Run("TestNodeGroupInstanceTypeIsInvalid", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "invalid-instance-type",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            Tags:           map[string]string{"key": "value"},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "InstanceType must be a valid EC2 instance type")
    })

    // Test case when a field is missing
    t.Run("TestNodeGroupRoleArnIsMissing", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            Tags:           map[string]string{"key": "value"},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "RoleArn must be specified")
    })

    // Test case when a field is empty
    t.Run("TestNodeGroupRoleArnIsEmpty", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "",
            Tags:           map[string]string{"key": "value"},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "RoleArn cannot be empty")
    })

    // Test case when a field is invalid
    t.Run("TestNodeGroupRoleArnIsInvalid", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "invalid-role-arn",
            Tags:           map[string]string{"key": "value"},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "RoleArn must be a valid IAM role ARN")
    })

    // Test case when a field is missing
    t.Run("TestNodeGroupTagsAreMissing", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Tags must be specified")
    })

    // Test case when a field is empty
    t.Run("TestNodeGroupTagsAreEmpty", func(t *testing.T) {
        nodeGroup := utils.NodeGroup{
            Name:            "my-node-group",
            InstanceType:    "t3.medium",
            DesiredCapacity: 2,
            MinSize:         1,
            MaxSize:         2,
            SubnetIds:       []string{"subnet-12345"},
            RoleArn:         "arn:aws:iam::123456789012:role/eks-node-group-role",
            Tags:           map[string]string{},
        }
        err := utils.ValidateClusterConfig(nodeGroup)
        require.Error(t, err)
        // require.Contains(t, err.Error(), "Tags cannot be empty")
    })
}
