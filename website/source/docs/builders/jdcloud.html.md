---
description: |
    The `jdcloud` Packer builder helps you to build instance images
    based on an existing image
layout: docs
page_title: 'JDCloud Image Builder'
sidebar_current: 'docs-builders-jdcloud'
---

# JDCloud Image Builder

Type: `jdcloud`

The `jdcloud` Packer builder helps you to build instance images
based on an existing image

## Configuration Reference

In order to build a JDCloud instance image, fullfill your configuration file. Necessary attributes
are given below: 

### Attributes:

- `type` (string) - This parameter tells which cloud-service-provider you are using, in our case, use 'jdcloud'
- `source_image_id` (string) - New image is generated based on an old one, specify the old one here. 
- `access_key` (string) - Declare your identity , help us to know who you are. Alternatively you can set them as env-variable:`export JDCLOUD_ACCESS_KEY=xxx`
- `secret_key` (string) - Same as `access_key`, write them this file, or set them as env-variable:`export JDCLOUD_SECRET_KEY=xxx`
- `region_id` (string) - Region of your instance, can be 'cn-north-1/cn-east-1' etc.
- `az` (string) - Exact availability zone of instance, 'cn-north-1c' for example
- `instance_name` (string) - Name your instance
- `image_name` (string) - Name the image you would like to create
- `subnet_id` (string) - An instance is supposed to exists in an subnet
- `password` (string) -  Password of your instance
- `ssh_wait_timeout`(string) - 'Timeout' is introduced in trying to connect(ssh) to your instance
- `provisioners/inline` (string) - Commands, were written in an array form, where each command is a string, e.g it can be`apt-get update`
- `communicator` (string) - Currently only `ssh` is supported. `winrm` will be added if required
- `ssh_username` (string) - Currently only `root` is supported 

## Examples

Here is a basic example for JDCloud.

``` json
{
  "variables": {
    "jdcloud_access_key": "<your-access-key>",
    "jdcloud_secret_key": "<your-secret-key>"
  },
  "builders": [
    {
      "type": "jdcloud",
      "source_image_id": "<some-image>",
      "access_key": "<your-access-key>",
      "secret_key": "<your-secret-key>",
      "region_id": "cn-north-1",
      "az": "cn-north-1c",
      "instance_name": "created_by_packer",
      "instance_type": "g.n2.medium",
      "image_name": "packerImage",
      "subnet_id": "<some-subnet>",
      "password": "DevOps2018",
      "communicator":"ssh",
      "ssh_username": "root",
      "ssh_password": "DevOps2018",
      "ssh_wait_timeout" : "60s"
    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "inline": [
        "sleep 30",
        "sudo apt-get install -y python-pip"
      ]
    }
  ]
}

```

[Find more examples](https://github.com/hashicorp/packer/tree/master/examples/jdcloud)