package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	vm "github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	vpc "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/models"
	"time"
)

// Plan
// Currently only these parameter are provided, prepare more
// Implement CleanUp function

type stepCreateJDCloudInstance struct {
	Az           string
	InstanceName string
	InstanceType string
	ImageId      string
	SubnetId     string
	InstanceId   string
}

func (s *stepCreateJDCloudInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now begin creating instances")

	generalConfig := state.Get("config").(*Config)
	vmClient := generalConfig.VmClient
	regionId := generalConfig.RegionId

	instanceSpec := vm.InstanceSpec{
		Az:           &s.Az,
		InstanceType: &s.InstanceType,
		ImageId:      &s.ImageId,
		Name:         s.InstanceName,
		PrimaryNetworkInterface: &vm.InstanceNetworkInterfaceAttachmentSpec{
			NetworkInterface: &vpc.NetworkInterfaceSpec{SubnetId: s.SubnetId, Az: &s.Az},
		},
	}
	req := apis.NewCreateInstancesRequest(regionId, &instanceSpec)
	resp, err := vmClient.CreateInstances(req)

	if err != nil || resp.Error.Code != 0 {
		err := fmt.Errorf("Error creating instance-Error-%s, Respond status-Code:%s,Message:%s", err,resp.Error.Code,resp.Error.Message)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	instanceId := resp.Result.InstanceIds[0]
	resultingStatus := waitForInstance(instanceId, regionId, vmClient, VMRunning)
	if resultingStatus != nil {
		err := fmt.Errorf("Error creating instance: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	reqInstance := apis.NewDescribeInstanceRequest(regionId, instanceId)
	respInstance, errInstance := vmClient.DescribeInstance(reqInstance)
	if errInstance != nil || respInstance.Error.Code != 0 {
		ui.Say(err.Error())
		return multistep.ActionHalt
	}

	s.InstanceId = respInstance.Result.Instance.InstanceId
	state.Put("instanceId", respInstance.Result.Instance.InstanceId)
	ui.Message(fmt.Sprintf("VM has been created, with instanceName:%s and instanceId:%s",
			   respInstance.Result.Instance.InstanceName,respInstance.Result.Instance.InstanceId))
	return multistep.ActionContinue
}

func waitForInstance(instanceId string, regionId string, vmClient *client.VmClient, expectedStatus string) error {
	currentTime := int(time.Now().Unix())
	req := apis.NewDescribeInstanceRequest(regionId, instanceId)
	connectFailedCount := 0
	for {
		time.Sleep(time.Second * 10)
		resp, err := vmClient.DescribeInstance(req)
		if resp.Result.Instance.Status == expectedStatus {
			return nil
		}
		if int(time.Now().Unix())-currentTime > Timeout {
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

func (s *stepCreateJDCloudInstance) Cleanup(state multistep.StateBag) {

}
