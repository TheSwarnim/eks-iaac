package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createOrUpdateCluster(ctx *pulumi.Context, clusterConfig utils.ClusterConfig) error {
	log.Printf("Creating EKS cluster: %s", clusterConfig.Name)

 	cluster, err := eks.NewCluster(ctx, clusterConfig.Name, &eks.ClusterArgs{
		Name:    pulumi.String(clusterConfig.Name),
		RoleArn: pulumi.String(clusterConfig.RoleArn),
		VpcConfig: &eks.ClusterVpcConfigArgs{
			PublicAccessCidrs: utils.ConvertToPulumiStringArray(clusterConfig.PublicAccessCidrs), // Convert []string to pulumi.StringArray
			SecurityGroupIds:  utils.ConvertToPulumiStringArray(clusterConfig.SecurityGroupIds),
			SubnetIds:         utils.ConvertToPulumiStringArray(clusterConfig.SubnetIds),
		},
		Version: pulumi.String(clusterConfig.Version),
		Tags: utils.ConvertToPulumiStringMap(clusterConfig.Tags), // Convert map[string]string to pulumi.StringMap
	})

	if err != nil {
		log.Printf("Failed to create EKS cluster: %s", clusterConfig.Name)
		return err
	}

	log.Printf("Successfully created EKS cluster: %s", clusterConfig.Name)

	err = createOrUpdateNodeGroups(ctx, clusterConfig, cluster)
	if err != nil {
		log.Printf("Failed to create node groups for cluster: %s", clusterConfig.Name)
		return err
	}

	return nil
}

func CreateOrUpdateClusters(ctx *pulumi.Context, clusterConfigs []utils.ClusterConfig) error {
	// Iterate over the clusterConfigs and create each cluster
	for _, clusterConfig := range clusterConfigs {
		log.Printf("Creating cluster: %s", clusterConfig.Name)

		// Use the CreateCluster function from src/components/cluster.go to create the cluster
		err := createOrUpdateCluster(ctx, clusterConfig)
		if err != nil {
			return err
		}

		log.Printf("Successfully created cluster: %s", clusterConfig.Name)
	}

	return nil
}