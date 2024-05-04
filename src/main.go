package main

import (
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/components"
	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new config object for the current Pulumi stack
		conf := config.New(ctx, "")

		// Read the root directory path from the Pulumi config
		rootDir := conf.Require("clusters-config-path")

		// Use the ReadClusterConfigs function from src/utils/readconfig.go to read the cluster configuration files
		clusterConfigs, err := utils.ReadClusterConfigs(rootDir)
		if err != nil {
			log.Fatalf("Failed to read cluster configs: %v", err)
		}

		// Iterate over the clusterConfigs and create each cluster
		for _, clusterConfig := range clusterConfigs {
			log.Printf("Creating cluster: %s", clusterConfig.Name)

			// Use the CreateCluster function from src/components/iam.go to create the IAM resources
			err := components.CreateOrUpdateCluster(ctx, clusterConfig)
			if err != nil {
				log.Fatalf("Failed to create cluster %s: %v", clusterConfig.Name, err)
			}

			log.Printf("Successfully created cluster: %s", clusterConfig.Name)
		}

		return nil
	})
}