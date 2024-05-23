package components

import (
	"fmt"
	"log"

	"github.com/dreamplug-tech/eks-iaac-2.0/src/utils"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam" // Add this import statement
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createClusterRole(ctx *pulumi.Context, roleName string, clusterConfig utils.ClusterConfig) (*iam.Role, error) {
    // Create the role
    role, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
		Name: pulumi.String(roleName),
        AssumeRolePolicy: pulumi.String(`{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "eks.amazonaws.com"
                    },
                    "Action": "sts:AssumeRole"
                }
            ]
        }`),
		Tags: utils.ConvertToPulumiStringMap(clusterConfig.Tags),
    })
    if err != nil {
        return nil, err
    }

    // Attach the AmazonEKSClusterPolicy managed policy
    _, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-policy", roleName), &iam.RolePolicyAttachmentArgs{
        Role:      role.Name,
        PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"),
		
    })
    if err != nil {
        return nil, err
    }

    return role, nil
}

// create a function that abstract the creation or getting of the role
func getOrCreateClusterRole(ctx *pulumi.Context, clusterConfig utils.ClusterConfig) (*iam.Role, error) {
	clusterRoleName := clusterConfig.Name + "-eks-cluster-role"
	if clusterConfig.RoleArn == "" {
		log.Println("RoleArn is empty, creating a new role")
		role, err := createClusterRole(ctx, clusterRoleName, clusterConfig)
		if err != nil {
			log.Fatalf("Failed to create role for cluster: %s", clusterConfig.Name)
			return nil, err
		}
		log.Println("Role creation successful")
		return role, nil
	} else {
		log.Println("RoleArn exists, using the existing role")
		role, err := iam.GetRole(ctx, clusterRoleName, pulumi.ID(clusterConfig.RoleArn), nil, nil)
		if err != nil {
			log.Fatalf("Failed to get the existing role %s for cluster: %s", clusterConfig.RoleArn, clusterConfig.Name)
			return nil, err
		}
		log.Println("Successfully got the existing role")
		return role, nil
	}
}

func createNodeGroupRole(ctx *pulumi.Context, roleName string, nodeGroupConfig utils.NodeGroupConfig) (*iam.Role, error) {
	role, err := iam.NewRole(ctx, roleName, &iam.RoleArgs{
		Name: pulumi.String(roleName),
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
				"Action": "sts:AssumeRole"
                }
            ]
        }`),
		Tags: utils.ConvertToPulumiStringMap(nodeGroupConfig.Tags),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-policy", roleName), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-policy-2", roleName), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"),
	})
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-policy-3", roleName), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func getOrCreateNodeGroupRole(ctx *pulumi.Context, nodeGroupConfig utils.NodeGroupConfig, clusterName string) (*iam.Role, error) {
	// the nodegroup role name should add the cluster name to make it unique
	nodeGroupRoleName := clusterName + "-" + nodeGroupConfig.Name + "-eks-nodegroup-role"
	// log.Println("NodeGroupRoleName: ", nodeGroupRoleName)
	if nodeGroupConfig.RoleArn == "" {
		log.Println("RoleArn is empty, creating a new role")
		role, err := createNodeGroupRole(ctx, nodeGroupRoleName, nodeGroupConfig)
		if err != nil {
			log.Fatalf("Failed to create role for nodegroup: %s", nodeGroupConfig.Name)
			return nil, err
		}
		log.Println("Role creation successful")
		return role, nil
	} else {
		log.Println("RoleArn exists, using the existing role")
		role, err := iam.GetRole(ctx, nodeGroupRoleName, pulumi.ID(nodeGroupConfig.RoleArn), nil, nil)
		if err != nil {
			log.Fatalf("Failed to get the existing role %s for nodegroup: %s", nodeGroupConfig.RoleArn, nodeGroupConfig.Name)
			return nil, err
		}
		log.Println("Successfully got the existing role")
		return role, nil
	}
}