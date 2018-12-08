package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"time"
)

// TODO
// Not being able to read public IP
// TODO You have to give it a public Ip

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
	Password     string
}

func (s *stepCreateJDCloudInstance) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {

	ui := state.Get("ui").(packer.Ui)
	ui.Say("Now begin creating instances")

	//generalConfig := state.Get("config").(Config)
	//vmClient := generalConfig.VmClient
	//regionId := generalConfig.RegionId
	//
	//instanceSpec := vm.InstanceSpec{
	//	Az:           &s.Az,
	//	InstanceType: &s.InstanceType,
	//	ImageId:      &s.ImageId,
	//	Name:         s.InstanceName,
	//	Password:     &s.Password,
	//	PrimaryNetworkInterface: &vm.InstanceNetworkInterfaceAttachmentSpec{
	//		NetworkInterface: &vpc.NetworkInterfaceSpec{SubnetId: s.SubnetId, Az: &s.Az},
	//	},
	//}
	//req := apis.NewCreateInstancesRequest(regionId, &instanceSpec)
	//resp, err := vmClient.CreateInstances(req)
	//
	//if err != nil || resp.Error.Code != 0 {
	//	err := fmt.Errorf("Error creating instance-Error-%s, Respond status-Code:%s,Message:%s", err, resp.Error.Code, resp.Error.Message)
	//	state.Put("error", err)
	//	ui.Error(err.Error())
	//	return multistep.ActionHalt
	//}
	//
	//instanceId := resp.Result.InstanceIds[0]
	//resultingStatus := waitForInstance(instanceId, regionId, vmClient, VMRunning)
	//if resultingStatus != nil {
	//	err := fmt.Errorf("Error creating instance: %s", err)
	//	state.Put("error", err)
	//	ui.Error(err.Error())
	//	return multistep.ActionHalt
	//}
	//
	//reqInstance := apis.NewDescribeInstanceRequest(regionId, instanceId)
	//respInstance, errInstance := vmClient.DescribeInstance(reqInstance)
	//if errInstance != nil || respInstance.Error.Code != 0 {
	//	ui.Say(err.Error())
	//	return multistep.ActionHalt
	//}
	//
	//s.InstanceId = respInstance.Result.Instance.InstanceId

	//state.Put("instanceId", respInstance.Result.Instance.InstanceId)
	//state.Put("publicIp", respInstance.Result.Instance.ElasticIpAddress)
	state.Put("instanceId", "i-t8bx190zeu")
	state.Put("publicIp", "116.196.95.69")

	//ui.Message(fmt.Sprintf("Instance has been created, Spec:{name:%s ,id:%s ,publicIp:%s}",
	//	respInstance.Result.Instance.InstanceName,
	//	respInstance.Result.Instance.InstanceId,
	//	respInstance.Result.Instance.ElasticIpAddress))
	ui.Message(fmt.Sprintf("Instance has been created, Spec:{name:%s ,id:%s ,publicIp:%s}",
		"haha",
		state.Get("instanceId").(string),
		state.Get("publicIp").(string)))
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
