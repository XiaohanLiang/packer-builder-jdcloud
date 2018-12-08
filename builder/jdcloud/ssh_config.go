package jdcloud

import (
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/template/interpolate"
)

type SSHConfig struct {
	Comm communicator.Config `mapstructure:",squash"`
}

func (c *SSHConfig) Prepare(ctx *interpolate.Context) []error {
	c.Comm.Prepare(ctx)
	return nil
}

func CommHost(state multistep.StateBag) (string, error) {
	ip := state.Get("publicIp").(string)
	return ip, nil
}
