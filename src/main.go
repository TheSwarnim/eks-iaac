package main

import (
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
			return err
		}

		// Create the clusters
		clusters, err := components.CreateOrUpdateClusters(ctx, clusterConfigs)
		if err != nil {
			return err
		}

		// Iterate over the clusters and read the nodegroup configurations
		for i := 0; i < len(clusterConfigs); i++ {
			// nodeGroupDirectory will be inside the cluster directory for each cluster
			nodeGroupDirectory := rootDir + "/" + clusterConfigs[i].Name + "/nodegroups"
			
			// Use the ReadNodeConfigs function from src/utils/readconfig.go to read the nodegroup configuration files
			nodeGroupConfigs, err := utils.ReadNodeConfigs(nodeGroupDirectory)
			if err != nil {
				return err
			}

			// Create the nodegroups for the current cluster
			err = components.CreateOrUpdateNodeGroups(ctx, nodeGroupConfigs, clusters[i], clusterConfigs[i].Name)
			if err != nil {
				return err
			}
		}

		return nil
	})
}