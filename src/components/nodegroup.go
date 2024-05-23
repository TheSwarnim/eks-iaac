package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createOrUpdateNodeGroup(ctx *pulumi.Context, nodeGroupConfig utils.NodeGroupConfig, cluster *eks.Cluster) (*eks.NodeGroup, error) {
	log.Printf("Creating or updating node group: %s", nodeGroupConfig.Name)

	nodeGroup, err := eks.NewNodeGroup(ctx, nodeGroupConfig.Name, &eks.NodeGroupArgs{
		ClusterName:   cluster.Name,
		NodeGroupName: pulumi.String(nodeGroupConfig.Name),
		NodeRoleArn:   pulumi.String(nodeGroupConfig.RoleArn),
		SubnetIds:     utils.ConvertToPulumiStringArray(nodeGroupConfig.NetworkConfiguration.SubnetIds),
		ScalingConfig: &eks.NodeGroupScalingConfigArgs{
			DesiredSize: pulumi.Int(nodeGroupConfig.ScalingConfiguration.DesiredCapacity),
			MinSize:     pulumi.Int(nodeGroupConfig.ScalingConfiguration.MinSize),
			MaxSize:     pulumi.Int(nodeGroupConfig.ScalingConfiguration.MaxSize),
		},
		InstanceTypes: pulumi.StringArray{pulumi.String(nodeGroupConfig.InstanceType)},
		Tags: 		utils.ConvertToPulumiStringMap(nodeGroupConfig.Tags),
	}, pulumi.DependsOn([]pulumi.Resource{cluster}))

	if err != nil {
		log.Printf("Failed to create or update node group: %s", nodeGroupConfig.Name)
		return nil, err
	}

	log.Printf("Successfully created or updated node group: %s", nodeGroupConfig.Name)

	return nodeGroup, nil
}

func createOrUpdateNodeGroups(ctx *pulumi.Context, nodeGroupConfigs []utils.NodeGroupConfig, cluster *eks.Cluster) error {
	for _, nodeGroupConfig := range nodeGroupConfigs {
		_, err := createOrUpdateNodeGroup(ctx, nodeGroupConfig, cluster)
		if err != nil {
			return err
		}
	}

	return nil
}