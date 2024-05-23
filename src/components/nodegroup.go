package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createOrUpdateNodeGroup(ctx *pulumi.Context, nodeGroupConfig utils.NodeGroupConfig, cluster *eks.Cluster, nodeGroupRole *iam.Role) (*eks.NodeGroup, error) {
	log.Printf("Creating or updating node group: %s", nodeGroupConfig.Name)

	nodeGroup, err := eks.NewNodeGroup(ctx, nodeGroupConfig.Name, &eks.NodeGroupArgs{
		ClusterName:   cluster.Name,
		NodeGroupName: pulumi.String(nodeGroupConfig.Name),
		NodeRoleArn:   nodeGroupRole.Arn,
		SubnetIds:     utils.ConvertToPulumiStringArray(nodeGroupConfig.NetworkConfiguration.SubnetIds),
		ScalingConfig: &eks.NodeGroupScalingConfigArgs{
			DesiredSize: pulumi.Int(nodeGroupConfig.ScalingConfiguration.DesiredCapacity),
			MinSize:     pulumi.Int(nodeGroupConfig.ScalingConfiguration.MinSize),
			MaxSize:     pulumi.Int(nodeGroupConfig.ScalingConfiguration.MaxSize),
		},
		InstanceTypes: utils.ConvertToPulumiStringArray(nodeGroupConfig.ComputeConfiguration.InstanceTypes),
		Tags: 		utils.ConvertToPulumiStringMap(nodeGroupConfig.Tags),
		Labels: 	utils.ConvertToPulumiStringMap(nodeGroupConfig.KubernetesLabels),
		Taints: 	utils.ConvertToPulumiTaintArray(nodeGroupConfig.KubernetesTaints),
		DiskSize: pulumi.Int(nodeGroupConfig.ComputeConfiguration.DiskSize),
		AmiType: pulumi.String(nodeGroupConfig.ComputeConfiguration.AmiType),
		RemoteAccess: &eks.NodeGroupRemoteAccessArgs{
			Ec2SshKey: pulumi.String(nodeGroupConfig.NetworkConfiguration.Ec2KeyPair),
		},
	}, pulumi.DependsOn([]pulumi.Resource{cluster}))

	if err != nil {
		log.Printf("Failed to create or update node group: %s", nodeGroupConfig.Name)
		return nil, err
	}

	log.Printf("Successfully created or updated node group: %s", nodeGroupConfig.Name)

	return nodeGroup, nil
}

func CreateOrUpdateNodeGroups(ctx *pulumi.Context, nodeGroupConfigs []utils.NodeGroupConfig, cluster *eks.Cluster) error {
	for _, nodeGroupConfig := range nodeGroupConfigs {
		log.Printf("Creating or updating node group: %s", nodeGroupConfig.Name)

		// check if roleArn is empty, if so, create a new role with suffix "-eks-nodegroup-role"
		nodeGroupRole, err := getOrCreateNodeGroupRole(ctx, nodeGroupConfig)
		if err != nil {
			return err
		}

		// Use the createOrUpdateNodeGroup function from src/components/nodegroup.go to create or update the nodegroup
		_, err = createOrUpdateNodeGroup(ctx, nodeGroupConfig, cluster, nodeGroupRole)
		if err != nil {
			return err
		}
	}

	return nil
}