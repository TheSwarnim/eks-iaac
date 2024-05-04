package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateOrUpdateCluster(ctx *pulumi.Context, clusterConfig utils.ClusterConfig) error {
	log.Printf("Creating EKS cluster: %s", clusterConfig.Name)

	_, err := eks.NewCluster(ctx, clusterConfig.Name, &eks.ClusterArgs{
		Name:    pulumi.String(clusterConfig.Name),
		RoleArn: pulumi.String(clusterConfig.RoleArn),
		VpcConfig: &eks.ClusterVpcConfigArgs{
			PublicAccessCidrs: utils.ConvertToPulumiStringArray(clusterConfig.PublicAccessCidrs), // Convert []string to pulumi.StringArray
			SecurityGroupIds:  utils.ConvertToPulumiStringArray(clusterConfig.SecurityGroupIds),
			SubnetIds:         utils.ConvertToPulumiStringArray(clusterConfig.SubnetIds),
		},
		Version: pulumi.String(clusterConfig.Version),
		Tags: func() pulumi.StringMap {
			result := make(pulumi.StringMap)
			for key, value := range clusterConfig.Tags {
				result[key] = pulumi.String(value)
			}
			return result
		}(),
	})

	if err != nil {
		log.Printf("Failed to create EKS cluster: %s", clusterConfig.Name)
		return err
	}

	log.Printf("Successfully created EKS cluster: %s", clusterConfig.Name)

	return nil
}
