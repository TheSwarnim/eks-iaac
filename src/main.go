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
		err = components.CreateOrUpdateClusters(ctx, clusterConfigs)
		if err != nil {
			return err
		}

		return nil
	})
}