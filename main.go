// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-dagger-cli/builder/scaffolding"
	scaffoldingData "github.com/hashicorp/packer-plugin-dagger-cli/datasource/scaffolding"
	scaffoldingPP "github.com/hashicorp/packer-plugin-dagger-cli/post-processor/scaffolding"
	daggercli "github.com/hashicorp/packer-plugin-dagger-cli/provisioner/dagger-cli"
	scaffoldingProv "github.com/hashicorp/packer-plugin-dagger-cli/provisioner/scaffolding"
	"github.com/hashicorp/packer-plugin-dagger-cli/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("my-builder", new(scaffolding.Builder))
	pps.RegisterProvisioner("my-provisioner", new(scaffoldingProv.Provisioner))
	pps.RegisterProvisioner("dagger-cli", new(daggercli.Provisioner))
	pps.RegisterPostProcessor("my-post-processor", new(scaffoldingPP.PostProcessor))
	pps.RegisterDatasource("my-datasource", new(scaffoldingData.Datasource))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
