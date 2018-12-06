package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	vm "github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
)

type stepStopJDCloudInstance struct{}

func (s *stepStopJDCloudInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now begin stopping this instance")

	generalConfig := state.Get("config").(*Config)
	instanceId := state.Get("instance").(*vm.Instance).InstanceId
	vmClient := generalConfig.AccessConfig.client
	regionId := generalConfig.AccessConfig.Region

	req := apis.NewStopInstanceRequest(regionId, instanceId)
	resp, err := vmClient.StopInstance(req)

	if err != nil || resp.Error.Code != 0 {
		err := fmt.Errorf("Error waiting for alicloud instance to stop: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	resultingStaus := waitForInstance(instanceId, regionId, vmClient, VMDeleted)
	if resultingStaus != nil {
		err := fmt.Errorf("Error waiting for alicloud instance to stop: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepStopJDCloudInstance) Cleanup(multistep.StateBag) {
	return
}
