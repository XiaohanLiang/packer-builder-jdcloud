package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	vm "github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	"time"
)

type stepCreateJDCloudImage struct {
	//image     *vm.Image
	ImageName string
}

func (s *stepCreateJDCloudImage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now begin stopping this instance")

	generalConfig := state.Get("config").(*Config)
	instanceId := state.Get("instance").(*vm.Instance).InstanceId
	vmClient := generalConfig.AccessConfig.client
	regionId := generalConfig.AccessConfig.Region

	req := apis.NewCreateImageRequest(regionId, instanceId, s.ImageName, "")
	resp, err := vmClient.CreateImage(req)

	if err != nil || resp.Error.Code != 0 {
		err := fmt.Errorf("Error creating image: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	imageId := resp.Result.ImageId
	resultingStatus := waitForImage(imageId, regionId, vmClient, ImageReady)

	if resultingStatus != nil {
		err := fmt.Errorf("Timeout waiting for image to be created: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("imageId", imageId)
	return multistep.ActionContinue
}

func waitForImage(imageId string, regionId string, vmClient *client.VmClient, expectedStatus string) error {
	currentTime := int(time.Now().Unix())
	req := apis.NewDescribeImageRequest(regionId, imageId)
	connectFailedCount := 0
	for {
		time.Sleep(time.Second * 10)
		resp, err := vmClient.DescribeImage(req)
		if resp.Result.Image.Status == expectedStatus {
			return nil
		}
		if int(time.Now().Unix())-currentTime > ImageTimeout {
			return fmt.Errorf("[ERROR] waitForInstance failed, timeout")
		}
		if err != nil {
			if connectFailedCount > Tolerance {
				return fmt.Errorf("[ERROR] waitForInstance, Tolerrance Exceeded failed %s ", err.Error())
			}
			connectFailedCount++
			continue
		} else {
			connectFailedCount = 0
		}
	}
	return nil
}

func (s *stepCreateJDCloudImage) Cleanup(state multistep.StateBag) {
	return
}
