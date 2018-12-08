package jdcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	vm "github.com/jdcloud-api/jdcloud-sdk-go/services/vm/models"
	vpc "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/models"
	vpcApis "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/apis"
	vpcClient "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/client"

	"time"
)

// Plan
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

	generalConfig := state.Get("config").(Config)
	vmClient := generalConfig.VmClient
	regionId := generalConfig.RegionId

	instanceSpec := vm.InstanceSpec{
		Az:           &s.Az,
		InstanceType: &s.InstanceType,
		ImageId:      &s.ImageId,
		Name:         s.InstanceName,
		Password:     &s.Password,
		PrimaryNetworkInterface: &vm.InstanceNetworkInterfaceAttachmentSpec{
			NetworkInterface: &vpc.NetworkInterfaceSpec{SubnetId: s.SubnetId, Az: &s.Az},
		},
	}
	req := apis.NewCreateInstancesRequest(regionId, &instanceSpec)
	resp, err := vmClient.CreateInstances(req)

	if err != nil || resp.Error.Code != 0 {
		err := fmt.Errorf("Error creating instance-Error-%s, Respond status-Code:%s,Message:%s", err, resp.Error.Code, resp.Error.Message)
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

	privateIpAddress :=  respInstance.Result.Instance.PrivateIpAddress
	networkInterfaceId := respInstance.Result.Instance.PrimaryNetworkInterface.NetworkInterface.NetworkInterfaceId
	s.InstanceId = respInstance.Result.Instance.InstanceId

	ui.Message("Creating publicIp...")
	publicIpId,errPublicIp := createElasticIp(state)
	if errPublicIp!=nil{
		ui.Say(errPublicIp.Error())
		return multistep.ActionHalt
	}

	ui.Message("Associating publicIp with this instance...")
	errAssociateIp := associatePublicIp(state,networkInterfaceId,publicIpId,privateIpAddress)
	if errAssociateIp!=nil{
		ui.Say(errAssociateIp.Error())
		return multistep.ActionHalt
	}

	reqPublicIpAddress := apis.NewDescribeInstanceRequest(regionId, instanceId)
	respPublicIpAddress, errPublicIpAddress := vmClient.DescribeInstance(reqPublicIpAddress)
	if errPublicIpAddress != nil || respPublicIpAddress.Error.Code != 0 {
		ui.Say(err.Error())
		return multistep.ActionHalt
	}

	state.Put("instanceId", respInstance.Result.Instance.InstanceId)
	state.Put("publicIp", respPublicIpAddress.Result.Instance.ElasticIpAddress)

	ui.Message(fmt.Sprintf("Instance has been created wuth spec:{name:%s ,id:%s ,publicIp:%s}",
		respInstance.Result.Instance.InstanceName,
		respInstance.Result.Instance.InstanceId,
		respPublicIpAddress.Result.Instance.ElasticIpAddress))
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

func createElasticIp(state multistep.StateBag)(string,error){

	generalConfig := state.Get("config").(Config)
	regionId := generalConfig.RegionId
	credential := core.NewCredentials(generalConfig.AccessKey, generalConfig.SecretKey)
	vpcclient := vpcClient.NewVpcClient(credential)

	req := vpcApis.NewCreateElasticIpsRequest(regionId,1,&vpc.ElasticIpSpec{
		BandwidthMbps:1,
		Provider:"bgp",
	})

	resp,err := vpcclient.CreateElasticIps(req)

	if err!=nil || resp.Error.Code!=0 {
		return "",fmt.Errorf("[ERROR] Failed in creating new publicIp, Error-%s, Response-Code:%s,Message:%s",err,resp.Error.Code,resp.Error.Message)
	}
	return resp.Result.ElasticIpIds[0],nil
}

func associatePublicIp(state multistep.StateBag,networkInterfaceId string,eipId string,privateIpAddress string) error {

	generalConfig := state.Get("config").(Config)
	regionId := generalConfig.RegionId
	credential := core.NewCredentials(generalConfig.AccessKey, generalConfig.SecretKey)
	vpcclient := vpcClient.NewVpcClient(credential)

	req := vpcApis.NewAssociateElasticIpRequest(regionId,networkInterfaceId)
	req.ElasticIpId = &eipId
	req.PrivateIpAddress = &privateIpAddress

	resp,err := vpcclient.AssociateElasticIp(req)

	if err!=nil || resp.Error.Code!=0 {
		return fmt.Errorf("[ERROR] Failed in associating publicIp, Error-%s, Response-Code:%s,Message:%s",err,resp.Error.Code,resp.Error.Message)
	}

	return nil
}