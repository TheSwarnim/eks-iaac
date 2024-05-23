package components

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam" // Add this import statement
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createOrUpdateCluster(ctx *pulumi.Context, clusterConfig utils.ClusterConfig, clusterRole *iam.Role) (*eks.Cluster, error) {
	log.Printf("Creating EKS cluster: %s", clusterConfig.Name)

 	cluster, err := eks.NewCluster(ctx, clusterConfig.Name, &eks.ClusterArgs{
		Name:    pulumi.String(clusterConfig.Name),
		RoleArn: clusterRole.Arn,
		KubernetesNetworkConfig: &eks.ClusterKubernetesNetworkConfigArgs{
			ServiceIpv4Cidr: pulumi.String(clusterConfig.ServiceIpv4Cidr),
		},
		VpcConfig: &eks.ClusterVpcConfigArgs{
			PublicAccessCidrs: utils.ConvertToPulumiStringArray(clusterConfig.PublicAccessCidrs), // Convert []string to pulumi.StringArray
			SecurityGroupIds:  utils.ConvertToPulumiStringArray(clusterConfig.SecurityGroupIds),
			SubnetIds:         utils.ConvertToPulumiStringArray(clusterConfig.SubnetIds),
		},
		Version: pulumi.String(clusterConfig.Version),
		Tags: utils.ConvertToPulumiStringMap(clusterConfig.Tags), // Convert map[string]string to pulumi.StringMap
        EnabledClusterLogTypes: pulumi.StringArray{pulumi.String("api"), pulumi.String("audit"), pulumi.String("authenticator"), pulumi.String("controllerManager"), pulumi.String("scheduler")},
	}, pulumi.DependsOn([]pulumi.Resource{clusterRole}))

	if err != nil {
		log.Printf("Failed to create EKS cluster: %s", clusterConfig.Name)
		return nil, err
	}

	log.Printf("Successfully created EKS cluster: %s", clusterConfig.Name)

	return cluster, nil
}

func CreateOrUpdateClusters(ctx *pulumi.Context, clusterConfigs []utils.ClusterConfig) ([]*eks.Cluster, error){
	var clusters []*eks.Cluster

	// Iterate over the clusterConfigs and create each cluster
	for _, clusterConfig := range clusterConfigs {
		log.Printf("Creating cluster: %s", clusterConfig.Name)

		// check if roleArn is empty, if so, create a new role with suffix "-eks-cluster-role"
		clusterRole, err := getOrCreateClusterRole(ctx, clusterConfig)
		if err != nil {
			return nil, err
		}
		
		// Use the CreateCluster function from src/components/cluster.go to create the cluster
		cluster, err := createOrUpdateCluster(ctx, clusterConfig, clusterRole)
		if err != nil {
			return nil, err
		}

		clusters = append(clusters, cluster)
	}

	return clusters, nil
}
