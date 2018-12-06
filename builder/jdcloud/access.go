package jdcloud

// TODO Currently abandoned
// Not quite sure if we are going to this place

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

// Invoke functions below and collect errors into an array form
func (c *JDCloudAccessConfig) Prepare(ctx *interpolate.Context) []error {

	errorArray := []error{}
	if err := c.Config(); err != nil {
		errorArray = append(errorArray, err)
	}

	if err := c.validateRegion(); err != nil {
		errorArray = append(errorArray, err)
	}

	if len(errorArray) != 0 {
		return errorArray
	}

	return nil
}

func (c *JDCloudAccessConfig) Config() error {

	if c.AccessKey == "" {
		c.AccessKey = os.Getenv("access_key")
	}

	if c.SecretKey == "" {
		c.SecretKey = os.Getenv("secret_key")
	}

	if c.AccessKey == "" || c.SecretKey == "" {
		return fmt.Errorf("[ERROR] Empty key pair, povide your keys or set as environment variable")
	}

	credential := core.NewCredentials(c.AccessKey, c.SecretKey)
	c.client = client.NewVmClient(credential)
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
