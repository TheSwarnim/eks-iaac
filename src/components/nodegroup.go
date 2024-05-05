package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createOrUpdateNodeGroup(ctx *pulumi.Context, nodeGroupConfig utils.NodeGroup, cluster *eks.Cluster) (*eks.NodeGroup, error) {
	log.Printf("Creating or updating node group: %s", nodeGroupConfig.Name)

	nodeGroup, err := eks.NewNodeGroup(ctx, nodeGroupConfig.Name, &eks.NodeGroupArgs{
		ClusterName:   cluster.Name,
		NodeGroupName: pulumi.String(nodeGroupConfig.Name),
		NodeRoleArn:   pulumi.String(nodeGroupConfig.RoleArn),
		SubnetIds:     getNodeGroupSubnetIds(nodeGroupConfig, cluster),
		ScalingConfig: &eks.NodeGroupScalingConfigArgs{
			DesiredSize: pulumi.Int(nodeGroupConfig.DesiredCapacity),
			MinSize:     pulumi.Int(nodeGroupConfig.MinSize),
			MaxSize:     pulumi.Int(nodeGroupConfig.MaxSize),
		},
		InstanceTypes: pulumi.StringArray{pulumi.String(nodeGroupConfig.InstanceType)},
		Tags: 		utils.ConvertToPulumiStringMap(nodeGroupConfig.Tags),
	})

	if err != nil {
		log.Printf("Failed to create or update node group: %s", nodeGroupConfig.Name)
		return nil, err
	}

	log.Printf("Successfully created or updated node group: %s", nodeGroupConfig.Name)

	return nodeGroup, nil
}

func createOrUpdateNodeGroups(ctx *pulumi.Context, clusterConfig utils.ClusterConfig, cluster *eks.Cluster) error {
	for _, nodeGroupConfig := range clusterConfig.NodeGroups {
		_, err := createOrUpdateNodeGroup(ctx, nodeGroupConfig, cluster)
		if err != nil {
			return err
		}
	}

	return nil
}

// if nodeGroupConfig.SubnetIds is not empty, use it as the subnetIds for the node group else use the cluster's subnetIds
func getNodeGroupSubnetIds(nodeGroupConfig utils.NodeGroup, cluster *eks.Cluster) pulumi.StringArrayInput {
	if len(nodeGroupConfig.SubnetIds) > 0 {
		return utils.ConvertToPulumiStringArray(nodeGroupConfig.SubnetIds)
	}

	return cluster.VpcConfig.SubnetIds()
}