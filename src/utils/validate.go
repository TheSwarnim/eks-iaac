package utils

import (
	"log"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// create a separate function for validation
func ValidateClusterConfig(config interface{}) error {
	validate := validator.New()
	err := validate.RegisterValidation("minfield", validateMaxSize)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("subnetid", validateSubnetID)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("instancetype", validateInstanceType)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("rolearn", validateRoleARN)
	if err != nil {
		return err
	}
	err = validate.RegisterValidation("securitygroupid", validateSecurityGroupID)
	if err != nil {
		return err
	}
	err = validate.Struct(config)
	if err != nil {
		return err
	}
	log.Println("Validation successful for config: ", config)
	return nil
}

// custom validation function to validate MaxSize field based on MinSize field: MaxSize >= MinSize
func validateMaxSize(fl validator.FieldLevel) bool {
    minFieldName := fl.Param()
    minSizeField := fl.Parent().FieldByName(minFieldName)
    if minSizeField.Kind() == reflect.Invalid {
        return false
    }
    minSize := minSizeField.Int()
    maxSize := fl.Field().Int()
    return maxSize >= minSize
}

// custom validation functions for SubnetId field
func validateSubnetID(fl validator.FieldLevel) bool {
    subnetID := fl.Field().String()
    // AWS subnet IDs start with "subnet-" followed by a 17-character hexadecimal string
    matched, _ := regexp.MatchString(`^subnet-[a-fA-F0-9]{17}$`, subnetID)
    return matched
}

// custom validation functions for InstanceType field
func validateInstanceType(fl validator.FieldLevel) bool {
    instanceType := fl.Field().String()
    // AWS instance types are in the format "t2.micro", "m5.large", etc.
    matched, _ := regexp.MatchString(`^[a-z]+\d\.[a-z]+$`, instanceType)
    return matched
}

// custom validation functions for RoleARN field
func validateRoleARN(fl validator.FieldLevel) bool {
    roleARN := fl.Field().String()
    // AWS IAM role ARNs are in the format "arn:aws:iam::123456789012:role/MyRole"
    matched, _ := regexp.MatchString(`^arn:aws:iam::\d{12}:role/.+$`, roleARN)
    return matched
}

// custom validation functions for security group ID field
func validateSecurityGroupID(fl validator.FieldLevel) bool {
	securityGroupID := fl.Field().String()
	// AWS security group IDs start with "sg-" followed by a 17-character hexadecimal string
	matched, _ := regexp.MatchString(`^sg-[a-fA-F0-9]{17}$`, securityGroupID)
	return matched
}