package jdcloud

import (
	"context"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

type stepPreValidate struct {
	JDCloudArtifactImageName string
	ForceDelete              bool
}

func (s *stepPreValidate) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now we begin validating")

	// PLAN Delete this step since there is no need to validate anything
	return multistep.ActionContinue

}
