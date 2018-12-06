package jdcloud

import (
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"log"
)

type Builder struct {
	config *Config
	runner multistep.Runner
}

// This function is invoked before builder begin running
// To make all data info prepared
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {

	config, errs := NewConfig(raws...)
	if errs != nil {
		return nil, errs
	}
	b.config = config
	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {

	ui.Say("Position-1 Begin building,and set up state bag")
	// Set up state bag
	state := new(multistep.BasicStateBag)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("config", b.config)

	ui.Say("Position-2 Begin building several steps")
	//ui.Say("Position-3")
	// Several step to execute different functions
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
		},
		&stepStopJDCloudInstance{},
		&stepCreateJDCloudImage{
			ImageName: b.config.ImageName,
		},
	}

	// Run these steps
	b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(state)

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

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
