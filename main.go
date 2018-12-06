package main

import (
	"github.com/XiaohanLiang/packer-builder-jdcloud/builder/jdcloud"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(jdcloud.Builder))
	server.Serve()
}
