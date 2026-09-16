// SPDX-FileCopyrightText: 2025 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"runtime/debug"

	"github.com/coder/terraform-provider-coderd/fwprovider"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
)

// terraformProviderModule is the Go module the coderd Terraform provider is
// embedded from. Its version in the build info is the version the embedded
// provider reports.
const terraformProviderModule = "github.com/coder/terraform-provider-coderd"

// NewFrameworkProvider returns a fresh, unconfigured instance of the coderd
// Terraform provider. The provider is linked into this binary and driven
// in-process through upjet's Terraform plugin framework client instead of being
// run as a plugin under the Terraform CLI.
func NewFrameworkProvider() tfprovider.Provider {
	return fwprovider.New(terraformProviderVersion())()
}

// terraformProviderVersion is the version of the embedded coderd Terraform
// provider, read from the module dependency recorded at build time.
func terraformProviderVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, dep := range bi.Deps {
		if dep.Path == terraformProviderModule {
			if dep.Replace != nil {
				return dep.Replace.Version
			}
			return dep.Version
		}
	}
	return "dev"
}
