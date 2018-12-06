package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
)

// Plan what does a config File look like
// 		parameter JDCloudSourceImageId is not used in this file

type stepCheckSourceImage struct {
	JDCloudSourceImageId string
}

func (s *stepCheckSourceImage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now begin validating source image")

	generalConfig := state.Get("config").(*Config)

	vmClient := generalConfig.VmClient
	regionId := generalConfig.RegionId
	sourceImageId := generalConfig.SourceImageId

	req := apis.NewDescribeImageRequest(regionId, sourceImageId)
	resp, err := vmClient.DescribeImage(req)

	if err != nil {
		err := fmt.Errorf("[ERROR] Validating source image failed : %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if resp.Error.Code != 0 || resp.Result.Image.ImageId == "" {
		err := fmt.Errorf("[ERROR] Source image not found, code:%s, message:%s", resp.Error.Code, resp.Error.Message)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf("Validation success, image found with id:%s ,name:%s ", sourceImageId, resp.Result.Image.Name))

	state.Put("source_image", &resp.Result.Image)
	return multistep.ActionContinue
}

func(stepCheckSourceImage) Cleanup(state multistep.StateBag) {}