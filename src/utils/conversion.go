package utils

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

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

func ConvertPulumiStringOutputToString(output pulumi.StringOutput) (string, error) {
	var result string
	output.ApplyT(func(value string) string {
		result = value
		return value
	})

	if result == "" {
		return "", errors.New("failed to convert pulumi.StringOutput to string")
	}
	return result, nil
}