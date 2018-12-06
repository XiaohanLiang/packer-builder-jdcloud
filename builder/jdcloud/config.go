package jdcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"github.com/mitchellh/mapstructure"
)

const (
	Timeout      = 300
	Tolerance    = 3
	VMRunning    = "running"
	VMDeleted    = "deleted"
	VMStopped    = "stopped"
	ImageTimeout = 300
	ImageReady   = "ready"
	BuilderID    = "JDCloud"
)

// Plan
// clear config info

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	SourceImageId       string `mapstructure:"source_image_id"`
	AccessKey           string `mapstructure:"access_key"`
	SecretKey           string `mapstructure:"secret_key"`
	RegionId            string `mapstructure:"region_id"`
	Az                  string `mapstructure:"az"`
	InstanceName        string `mapstructure:"instance_name"`
	InstanceType        string `mapstructure:"instance_type"`
	ImageName           string `mapstructure:"image_name"`
	SubnetId            string `mapstructure:"subnet_id"`
	VmClient            *client.VmClient
	ctx                 interpolate.Context
}

// generate Config variable
func NewConfig(raws ...interface{}) (*Config, error) {

	c := Config{}
	md := mapstructure.Metadata{}

	err := config.Decode(&c, &config.DecodeOpts{
		Metadata:    &md,
		Interpolate: true,
	}, raws...)

	if err != nil {
		return nil, err
	}

	// Accumulate Errors
	errArray := &packer.MultiError{}

	if len(errArray.Errors) > 0 {
		return nil, errArray
	}

	// Generate client
	credential := core.NewCredentials(c.AccessKey, c.SecretKey)
	c.VmClient = client.NewVmClient(credential)

	return &c, nil

}
