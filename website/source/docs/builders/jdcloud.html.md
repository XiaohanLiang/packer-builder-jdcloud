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

### Required-Parameters:

- `type` (string) - This parameter tells which cloud-service-provider you are using, in our case, use 'jdcloud'
- `image_id` (string) - New image is generated based on an old one, specify the base-image-id here. 
- `access_key` (string) - Declare your identity , help us to know who you are. Alternatively you can set them as env-variable:`export JDCLOUD_ACCESS_KEY=xxx`
- `secret_key` (string) - Same as access_key, write them this file, or set them as env-variable:`export JDCLOUD_SECRET_KEY=xxx`
- `region_id` (string) - Region of your instance, can be 'cn-north-1/cn-east-1' etc.
- `az` (string) - Exact availability zone of instance, 'cn-north-1c' for example
- `instance_name` (string) - Name your instance
- `instance_type` (string) - Class of your expected instance
- `image_name` (string) - Name the image you would like to create
- `ssh_username` (string) - Currently only `root` is supported 
- `communicator` (string) - Currently only `ssh` is supported. `winrm` will be added if required
- `provisioners/inline` (string) - Commands, were written in an array form, where each command is a string, e.g it can be`apt-get update`

### Optional-Parameters

- `subnet_id` (string) - An instance is supposed to exists in an subnet, if not specified , we will create new one for you
- `ssh_wait_timeout`(string) - 'Timeout' is introduced in trying to connect(ssh) to your instance
- Credentials: This space specify which way you would like us to login, they have to be one of the following options:
    - `password` (string) -  Login with password
    - `ssh_private_key_file` + `ssh_keypair_name` - Login with an existing keypair, you have to give the path to your private key and its key name\
    - `temporary_key_pair_name` - We will create a new key for you, and use them as login credentisl


## Examples

Here is a basic example for JDCloud.

``` json
{
  "builders": [
    {
      "type": "jdcloud",
      "image_id": "img-h8ly274zg9",
      "access_key": "E1AD46FF7F16620B339CF1C2C21AFA3D",
      "secret_key": "B527396D7B6AF3685A7CDD809714DAA7",
      "region_id": "cn-north-1",
      "az": "cn-north-1c",
      "instance_name": "created_by_packer",
      "instance_type": "g.n2.medium",
      "ssh_password":"/Users/mac/.ssh/id_rsa",
      "ssh_keypair_name":"created_by_xiaohan",
      "image_name": "packerImage",
      "subnet_id": "subnet-jo6e38sdli",
      "communicator":"ssh",
      "ssh_username": "root",
      "ssh_timeout" : "60s"
    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "inline": [
        "sleep 3",
        "echo 123",
        "pwd"
      ]
    }
  ]
}


```

[Find more examples](https://github.com/hashicorp/packer/tree/master/examples/jdcloud)