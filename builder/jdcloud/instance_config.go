package jdcloud

import (
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type JDCloudInstanceSpecConfig struct {
	SourceImageId   string              `mapstructure:"source_image_id"`
	InstanceName    string              `mapstructure:"instance_name"`
	InstanceType    string              `mapstructure:"instance_type"`
	ImageId         string              `mapstructure:"image_id"`
	ImageName       string              `mapstructure:"image_name"`
	Password        string              `mapstructure:"password"`
	SubnetId        string              `mapstructure:"subnet_id"`
	Comm            communicator.Config `mapstructure:",squash"`
	Communicator    string              `mapstructure:"communicator"`
	InstanceId      string
	ArtifactId      string
	PublicIpAddress string
}

func (jd *JDCloudInstanceSpecConfig) Prepare(ctx *interpolate.Context) []error {

	// validate images non nil, etc...

	return jd.Comm.Prepare(ctx)
}
