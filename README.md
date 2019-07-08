# Packer-Builder-JDCloud

## About Packer

- **TL;DR** Packer helps you to build VM images
- Packer official website: [Click here](www.packer.io)

## About this branch

We are trying to build our plugin for Hashicorp-Packer. This is a developer branch of packer-builder-jdcloud plugin.

## How can I access to this plugin

- From binary:
    - Download binary in the release page, then 
    - Follow [this instruction](https://www.packer.io/docs/extending/plugins.html#installing-plugins) to use unreleased version
    - Start [debugging](https://www.packer.io/docs/other/debugging.html)
- From source code(recommended only when familiar with golang)
```bash
cd packer-builder-jdcloud
make build
```