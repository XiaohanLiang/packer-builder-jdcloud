package jdcloud

import (
	"fmt"
	"context"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
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

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	SSHConfig           `mapstructure:",squash"`
	SourceImageId       string `mapstructure:"source_image_id"`
	AccessKey           string `mapstructure:"access_key"`
	SecretKey           string `mapstructure:"secret_key"`
	RegionId            string `mapstructure:"region_id"`
	Az                  string `mapstructure:"az"`
	InstanceName        string `mapstructure:"instance_name"`
	InstanceType        string `mapstructure:"instance_type"`
	ImageName           string `mapstructure:"image_name"`
	Password            string `mapstructure:"password"`
	SubnetId            string `mapstructure:"subnet_id"`

	Communicator string `mapstructure:"communicator"`
	SSH_Username string `mapstructure:"ssh_username"`
	SSH_Password string `mapstructure:"ssh_password"`
	SSH_Timeout  string `mapstructure:"ssh_wait_timeout"`

	VmClient *client.VmClient
	ctx      interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

// This function is invoked before builder begin running
// To make all data info prepared
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"boot_command",
			},
		},
	}, raws...)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed in decoding JSON->mapstructure")
	}

	errs := &packer.MultiError{}
	errs = packer.MultiErrorAppend(errs, b.config.SSHConfig.Prepare(&b.config.ctx)...)
	if errs != nil && len(errs.Errors) != 0 {
		return nil, errs
	}

	packer.LogSecretFilter.Set(b.config.AccessKey, b.config.SecretKey)
	credential := core.NewCredentials(b.config.AccessKey, b.config.SecretKey)
	b.config.VmClient = client.NewVmClient(credential)

	return nil, nil
}

func (b *Builder) Run(ctx context.Context,ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {

	ui.Say("Position-1 Begin building,and set up state bag")
	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("config", b.config)

	ui.Say("Position-2 Begin building several steps")

	steps := []multistep.Step{

		&stepCheckSourceImage{
			JDCloudSourceImageId: b.config.SourceImageId,
		},

		&stepCreateJDCloudInstance{
			Az:           b.config.Az,
			InstanceName: b.config.InstanceName,
			InstanceType: b.config.InstanceType,
			ImageId:      b.config.SourceImageId,
			SubnetId:     b.config.SubnetId,
			Password:     b.config.Password,
		},

		//Config the communicator-SSH
		&communicator.StepConnect{
			Config:    &b.config.SSHConfig.Comm,
			Host:      CommHost,
			SSHConfig: b.config.SSHConfig.Comm.SSHConfigFunc(),
		},

		// Now we begin provisioning process
		&common.StepProvision{},

		&stepStopJDCloudInstance{},
		&stepCreateJDCloudImage{
			ImageName: b.config.ImageName,
		},
	}

	// Run these steps
	b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(ctx,state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	artifact := &Artifact{
		ImageId:  state.Get("imageId").(string),
		RegionID: b.config.RegionId,
	}
	return artifact, nil
}

//func (b *Builder) Cancel() {
//	if b.runner != nil {
//		log.Println("Cancelling the step runner...")
//		b.runner.Cancel()
//	}
//}
