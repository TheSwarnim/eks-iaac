package utils

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

// Helper function to convert []string to pulumi.StringArray
func ConvertToPulumiStringArray(items []string) pulumi.StringArray {
	result := make(pulumi.StringArray, len(items))
	for i, item := range items {
		result[i] = pulumi.String(item)
	}
	return result
}

func ConvertToPulumiStringMap(items map[string]string) pulumi.StringMap {
	result := make(pulumi.StringMap)
	for key, value := range items {
		result[key] = pulumi.String(value)
	}
	return result
}