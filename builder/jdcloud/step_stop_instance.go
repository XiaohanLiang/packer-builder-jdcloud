package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
)

type stepStopJDCloudInstance struct {
}

func (s *stepStopJDCloudInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Process - stepStopJDCloudInstance")

	generalConfig := state.Get("config").(Config)
	instanceId := state.Get("instanceId").(string)
	vmClient := generalConfig.VmClient
	regionId := generalConfig.RegionId

	req := apis.NewStopInstanceRequest(regionId, instanceId)
	resp, err := vmClient.StopInstance(req)
	if err != nil || resp.Error.Code != 0 {
		err := fmt.Errorf("[ERROR] Trying to stop this vm: Error-%s ,Resp-code:%s, message:%s", err, resp.Error.Code, resp.Error.Message)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message("Trying to stop this VM...")
	stoppingStatus := waitForInstance(instanceId, regionId, vmClient, VMStopped)
	if stoppingStatus != nil {
		err := fmt.Errorf("[ERROR] Waiting for JDCloud instance to stop: err:%s", stoppingStatus)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("VM with id %s has been stopped", instanceId))
	return multistep.ActionContinue
}

func (s *stepStopJDCloudInstance) Cleanup(multistep.StateBag) {
	return
}
