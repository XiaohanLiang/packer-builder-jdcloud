package jdcloud

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/mitchellh/mapstructure"
)

const (
	Timeout      = 300
	Tolerance    = 3
	VMRunning    = "running"
	VMDeleted    = ""
	ImageTimeout = 300
	ImageReady   = "ready"
	BuilderID    = "JDCloud"
)

// Plan
// clear config info

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ConfigFile          string   `mapstructure:"config_file"`
	OutputDir           string   `mapstructure:"output_directory"`
	ContainerName       string   `mapstructure:"container_name"`
	CommandWrapper      string   `mapstructure:"command_wrapper"`
	RawInitTimeout      string   `mapstructure:"init_timeout"`
	CreateOptions       []string `mapstructure:"create_options"`
	StartOptions        []string `mapstructure:"start_options"`
	AttachOptions       []string `mapstructure:"attach_options"`
	Name                string   `mapstructure:"template_name"`
	Parameters          []string `mapstructure:"template_parameters"`
	EnvVars             []string `mapstructure:"template_environment_vars"`
	TargetRunlevel      int      `mapstructure:"target_runlevel"`
	InitTimeout         time.Duration

	SourceImageId string `mapstructure:"source_image_id"`
	AccessConfig  JDCloudAccessConfig
	Az            string
	InstanceName  string
	InstanceType  string
	ImageName     string
	SubnetId      string
	ctx           interpolate.Context
}

// generate Config variable
func NewConfig(raws ...interface{}) (*Config, error) {

	c := Config{}
	md := mapstructure.Metadata{}

	err := config.Decode(c, &config.DecodeOpts{
		Metadata:    &md,
		Interpolate: true,
	}, raws...)

	if err != nil {
		return nil, err
	}

	// Accumulate Errors
	errArray := &packer.MultiError{}

	if c.OutputDir == "" {
		c.OutputDir = "output-" + c.PackerBuildName
	}

	if c.ContainerName == "" {
		c.ContainerName = "packer-" + c.PackerBuildName
	}

	if c.TargetRunlevel == 0 {
		c.TargetRunlevel = 3
	}

	if c.CommandWrapper == "" {
		c.CommandWrapper = "{{.Command}}"
	}

	if c.RawInitTimeout == "" {
		c.RawInitTimeout = "20s"
	}

	c.InitTimeout, err = time.ParseDuration(c.RawInitTimeout)
	if err != nil {
		errArray = packer.MultiErrorAppend(errArray, fmt.Errorf("Failed parsing init_timeout: %s", err))
	}

	if _, err := os.Stat(c.ConfigFile); os.IsNotExist(err) {
		errArray = packer.MultiErrorAppend(errArray, fmt.Errorf("LXC Config file appears to be missing: %s", c.ConfigFile))
	}

	if len(errArray.Errors) > 0 {
		return nil, errArray
	}

	return &c, nil

}
