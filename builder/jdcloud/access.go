package jdcloud

import (
	"fmt"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"os"
)

type JDCloudAccessConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	client    *client.VmClient
}

// Build a client for JDCloud client
//func (config *JDCloudAccessConfig) Client() (*client.VmClient, error) {
//vmClient := client.NewVmClient(config.Credential)
//return vmClient, nil
//}

// Invoke functions below and collect errors into an array form
func (config *JDCloudAccessConfig) Prepare(ctx *interpolate.Context) []error {

	errorArray := []error{}
	if err := config.Config(); err != nil {
		errorArray = append(errorArray, err)
	}

	if err := config.validateRegion(); err != nil {
		errorArray = append(errorArray, err)
	}

	if len(errorArray) != 0 {
		return errorArray
	}

	return nil
}

func (config *JDCloudAccessConfig) Config() error {

	if config.AccessKey == "" {
		config.AccessKey = os.Getenv("access_key")
	}

	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("secret_key")
	}

	if config.AccessKey == "" || config.SecretKey == "" {
		return fmt.Errorf("[ERROR] Empty key pair, povide your keys or set as environment variable")
	}

	credential := core.NewCredentials(config.AccessKey, config.SecretKey)
	config.client = client.NewVmClient(credential)
	return nil
}

func (config *JDCloudAccessConfig) validateRegion() error {
	regionArray := []string{"cn-north-1", "cn-south-1", "cn-east-1", "cn-east-2"}
	for _, item := range regionArray {
		if item == config.Region {
			return nil
		}
	}
	return fmt.Errorf("[ERROR] Invalid Region detected: %s", config.Region)
}
