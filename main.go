package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/XiaohanLiang/packer-builder-jdcloud/builder/jdcloud"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(jdcloud.Builder))
	server.Serve()
}